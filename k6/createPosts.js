import {check} from 'k6'
import http from 'k6/http';
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.3/index.js";
import { SharedArray } from 'k6/data';

const config = JSON.parse(open('../config/config.json'));
const creds = JSON.parse(open('../temp_store.json'));

export var options = {};
const {RPS, Executor, Duration, Rate, TimeUnit, VirtualUserCount, BatchSize} = config.LoadTestConfiguration;
const {MaxWordsCount, MaxWordLength} = config.PostsConfiguration;

// Temporary code
const shared = new SharedArray('demo array', () => {
    // const creds = JSON.parse(open('../temp_store.json'));
    return [0, 0];
})

if (RPS) {
    options = {
        discardResponseBodies: false,
        scenarios: {
            contacts: {
                executor: Executor,
                duration: Duration,
                rate: Rate,
                timeUnit: TimeUnit,
                preAllocatedVUs: VirtualUserCount,
            },
        }
    }
} else {
    options = {
        vus: VirtualUserCount,
        duration: Duration,
    }
}

export function setup() {
	if (MaxWordsCount <= 0) {
        console.error("Error in validating the posts configuration:", "max word count should be greater than 0");
		return;
	}

	if (MaxWordLength <= 0) {
        console.error("Error in validating the posts configuration:", "max word length should be greater than 0");
		return;
	}
    
    if (BatchSize <= 0 || BatchSize > 20) {
        console.error("Error in validating the batch size:", "batch size should be between 1 and 20");
        return;
    }
}

function getRandomMessage(wordsCount, wordLength) {
	let message = '';
    const characterSet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
	let words = 0;
    wordsCount = Math.floor(Math.random() * wordsCount) + 1;
    while (words < wordsCount) {
        let count = 0;
        wordLength = Math.floor(Math.random() * wordLength) + 1;
        while (count < wordLength) {
            message += characterSet.charAt(Math.floor(Math.random() * characterSet.length));
            count++;
        }

        message += ' ';
        words++;
    }

	return message;
}

function getRandomToken() {
    let tokens = [];
    creds.Users.map((u) => tokens.push(u.token));
    return tokens[(Math.floor(Math.random() * tokens.length))];
}

function getRandomChannel() {
    let channels = []
    if (creds.DM) {
        channels.push(creds.DM.id);
    }

    if (creds.GM) {
        channels.push(creds.GM.id);
    }

    creds.Channels.map((c) => channels.push(c));
    return channels[(Math.floor(Math.random() * channels.length))];
}

// Temporary code
export function handleSummary(data) {
    data.metrics['no_of_posts_created'] = {
        type: 'rate',
        contains: "default",
        values: {
            fails: shared[0],
            rate: 0,
            passes: shared[1],
        },
    }
    return {
      'stdout': textSummary(data, {indent: '    ', enableColors: true})
    };
}

export default function() {
    const requests = [];
    for(let i = 0; i < BatchSize; i++) {
        const channel = getRandomChannel();
        let url;
        if (typeof channel === "string") {
            const chatId = channel;
            url = `/chats/${chatId}/messages`;
        } else {
            const {id: channelId, team_id: teamId} = channel;
            url = `/teams/${teamId}/channels/${channelId}/messages`;
        }

        requests.push({
            url,
            method: "POST",
            id: i,
            body: {
                body: {
                    content: getRandomMessage(MaxWordsCount, MaxWordLength)
                }
            },
            headers: {
                "Content-Type": "application/json"
            }
        })
    }

    const payload = JSON.stringify({requests});
    const headers = {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${getRandomToken()}`,
    }

    const resp = http.post("https://graph.microsoft.com/v1.0/$batch", payload, {headers});
    check(resp, {
        'Response status is 200': (r) => resp.status === 200,
        'Response Content-Type header': (r) => {
            if(resp.headers['Content-Type']) {
                return resp.headers['Content-Type'].includes('application/json')
            }
        }
    });

    // TODO: Temporary code, fix or remove
    const responses = resp.json("responses");
    if(responses) {
        for(let response of responses) {
            if(response.status === 201) {
                shared[0] = shared[0]+1;
            } else {
                shared[1] = shared[1]+1;
            }
        }
    }
    
    
    let countSucc = shared[1], countFail = shared[0];
    console.info(`No. of posts created: ✔ ${countSucc}  ⨯ ${countFail}`);
}
