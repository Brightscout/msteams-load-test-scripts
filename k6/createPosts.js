import {check} from 'k6'
import http from 'k6/http';

const config = JSON.parse(open('../config/config.json'));
const creds = JSON.parse(open('../temp_store.json'));

export var options = {};
if (config.LoadTestConfiguration.RPS) {
    options = {
        discardResponseBodies: true,
        scenarios: {
            contacts: {
                executor: config.LoadTestConfiguration.Executor,
                duration: config.LoadTestConfiguration.Duration,
                rate: config.LoadTestConfiguration.Rate,
                timeUnit: config.LoadTestConfiguration.TimeUnit,
                preAllocatedVUs: config.LoadTestConfiguration.VirtualUserCount,
            },
        }
    }
} else {
    options = {
        vus: config.LoadTestConfiguration.VirtualUserCount,
        duration: config.LoadTestConfiguration.Duration,
    }
}

export function setup() {
	if (config.PostsConfiguration.MaxWordsCount <= 0) {
        console.error("Error in validating the posts configuration:", "max word count should be greater than 0");
		return;
	}

	if (config.PostsConfiguration.MaxWordLength <= 0) {
        console.error("Error in validating the posts configuration:", "max word length should be greater than 0");
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
    const payload = JSON.stringify({
        body: {
          content: getRandomMessage(config.PostsConfiguration.MaxWordsCount, config.PostsConfiguration.MaxWordLength)
        }
      });

    const headers = {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${getRandomToken()}`,
    }

    const channel = getRandomChannel();
    let resp;
    if (typeof channel === "string") {
        const chatId = channel;
        resp = http.post(`https://graph.microsoft.com/v1.0/chats/${chatId}/messages`, payload, {headers});
    } else {
        const {id, team_id} = channel;
        resp = http.post(`https://graph.microsoft.com/v1.0/teams/${team_id}/channels/${id}/messages`, payload, {headers});
    }


    check(resp, {
        'Post status is 201': (r) => resp.status === 201,
        'Post Content-Type header': (r) => resp.headers['Content-Type'].includes('application/json'),
      });
}
