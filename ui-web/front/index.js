import { render, $ } from "./util.mjs";
import collection from "./collection.mjs";

$main(collection);

function $main(state) {
	render([
		$("aside.tree",
			$("button.save", "Save"),
			$("button.save", "Run all"),
			[state].map($treeItem),
		),

		$("div.content",
			{ $: "textarea", name: "head", rows: 10, placeholder: "head" },
			{ $: "textarea", name: "body", rows: 10, placeholder: "body" },
			{ $: "textarea", name: "test", rows: 10, placeholder: "test" },
			$("button.send", "Send"),
			{ $: "code.status", title: "HTTP Status", t: "200 OK" },
			$("output.result-head", "result-head"),
			$("output.result-body", "result-body"),
			$("output.result-test", "result-test"),
		),
	])
}

function $treeItem(col) {
	if (!col) return;
	return $("div.tree-item",
		$("div.tree-name",
			col.name,
			" ", $("button.execute", "?>"),
			" ", $("button.execute", "X>"),
			" ", $("button.execute", "V>"),
			" ", $("button.execute", "C+"),
			" ", $("button.execute", "R+"),
			" ", $("button.execute", "-"),
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
		],
		onclick() { $$content(request) },
	};
}

function $$content(request) {
	render([
		{ $: "textarea", name: "head", rows: 10, placeholder: "head", t: headContent(request) },
		{ $: "textarea", name: "body", rows: 10, placeholder: "body" },
		{ $: "textarea", name: "test", rows: 10, placeholder: "test" },
		$("button.send", "Send"),
		{ $: "code.status", title: "HTTP Status", t: "200 OK +++++++++++" },
		$("output.result-head", "result-head"),
		$("output.result-body", "result-body"),
		$("output.result-test", "result-test"),
	], document.querySelector("div.content"))
}

function headContent(request) {
	return [
		request.method + " " + request.URL,
		// ...(request.header || [])
	].join("\n");
}