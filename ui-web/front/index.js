import { render, $ } from "./util.mjs";

$$main(fetch("/api/collection.json"));

async function $$main(request) {

	render("loading ...");

	const collection = await request.then(rep => rep.json());

	render([
		$("aside.tree",
			$b("button.bigAct", "Save", () => {
				fetch("/api/collection.json", {
					method: "PUT",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify(collection)
				}).then(({ ok }) => alert("saving: " + ok))
			}),
			$b("button.bigAct", "Run all", () => {
				$$main(fetch("/api/exe", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify(collection)
				}))
			}),
			[collection].map($treeItem),
		),
		$("div#content"),
	]);

	$$content(collection.children[0].requests[0]);
}

// Create button componenent with click handler.
function $b(n, t, onclick) {
	return { $: n, t, onclick };
}

/* TREE ZONE */

function $treeItem(col) {
	if (!col) return;
	return $("div.tree-item",
		$("div.tree-name",
			col.name,
			// " ", $b("button.execute", "C+", () => {
			// 	// col.request = (col.request ?? [] );
			// 	const name = prompt("name of the collection");
			// 	console.log('name:', name);
			// 	col.request ??= [];
			// 	col.request.push([{ name }]);
			// }),
			// " ", $b("button.execute", "C+"),
			// " ", $b("button.execute", "R+"),
			// " ", $b("button.execute", "-"),
		),
		(col.children || []).map($treeItem),
		(col.requests || []).map($treeRequest),
	);
}

function $treeRequest(request) {
	return {
		$: "div.tree-item",
		c: [
			$("div.tree-name", request.name),
			$("span.status",
				!request.response ? "[?]" : !request.response.testFails ? "[V]" : "[X]"
			),
		],
		onclick() { $$content(request) },
	};
}

/* CONTENT ZONE */


function $$content(request) {
	Object.defineProperty(request, "head", {
		configurable: true,
		enumerable: true,
		get() { return this.method + " " + this.url + "\n" + formatHeader(this.header) },
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
		}
	});


	const { response } = request;
	render([
		$$contentTextarea(request, "head"),
		$$contentTextarea(request, "body"),
		$$contentTextarea(request, "test"),
		response && [
			$("output.result-head", response.statusCode + "\n" + formatHeader(response.header)),
			$("output.result-body", atob(response.body || "")),
			$("output.result-test", (response.testFails || []).map(fail => $("div.fail", fail))),
		],
	], "content");
}

function $$contentTextarea(request, name) {
	return {
		$: "textarea", rows: 10,
		name, placeholder: name,
		value: request[name] || "",
		oninput({ target: { value } }) { request[name] = value; },
	};
}

function headContent(request) {
	return [
		request.method + " " + request.URL,
		formatHeader(request.header)
	].join("\n");
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