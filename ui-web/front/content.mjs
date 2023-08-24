import { render, $ } from "./util.mjs";

export function defineHeadPropertie(request) {
	Object.defineProperty(request, "head", {
		configurable: true,
		enumerable: true,
		get() { return (this.method ?? "GET") + " " + (this.url ?? "http://localhost:8080/") + "\n" + formatHeader(this.header) },
		set(value) {
			const [first, ...headers] = value.split("\n");
			let [_, method, url] = /(\w+)\s+(.*)/.exec(first);
			this.method = method;
			this.url = url;
			this.header = {};
			for (const line of headers) {
				const [key, value = ""] = line.split(':', 2).map(s => s.trim()),
					values = this.header[key] ?? [];
				if (!key || !value) continue
				values.push(value);
				this.header[key] = values;
			}
		},
	});

	Object.defineProperty(request, "auth", {
		configurable: true,
		enumerable: true,
		get() { return JSON.stringify(request.authorization ?? {}, null, "\t") },
		set(value) { this.authorization = JSON.parse(value); },
	})
}

// Print the content of one request
export function $$content(request) {
	const { response } = request;
	render([
		$contentTextarea(request, "head"),
		$contentTextarea(request, "auth"),
		$contentTextarea(request, "body"),
		$contentTextarea(request, "test"),
		response && [
			$("output.result-head", response.statusCode + "\n" + formatHeader(response.header)),
			$("output.result-body", atob(response.body || "")),
			$("output.result-test", (response.testFails || []).map(fail => $("div.fail", fail))),
		],
	], "content");
}

function $contentTextarea(request, name) {
	return {
		$: "textarea", rows: 10,
		name, placeholder: name,
		value: request[name] || "",
		oninput({ target: { value } }) { request[name] = value; },
	};
}

function $auth(request) {
	return $("div", "auth ...");
}

function formatHeader(header = {}) {
	let out = "";
	for (const key in header) {
		for (const value of header[key]) {
			out += key + ": " + value + "\n";
		}
	}
	return out;
}