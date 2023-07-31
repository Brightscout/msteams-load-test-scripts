import {check} from 'k6'
import http from 'k6/http';

const config = JSON.parse(open('../config/config.json'));
const creds = JSON.parse(open('../temp_store.json'));

export var options = {};
const {RPS, Executor, Duration, Rate, TimeUnit, VirtualUserCount, BatchSize} = config.LoadTestConfiguration;
const {MaxWordsCount, MaxWordLength} = config.PostsConfiguration;

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

    let countSucc = 0, countFail = 0;
    const responses = resp.json("responses");
    if(responses) {
        for(let response of responses) {
            if(response.status === 201) {
                countSucc++;
            } else {
                countFail++;
            }
        }
    }
    
    console.info(`No. of posts created: ✔ ${countSucc}  ⨯ ${countFail}`);
}
