import { render, $, $b } from "./util.mjs";
import { $$content, defineHeadPropertie } from "./content.mjs";

const PARENT = Symbol("parent");

$$main(fetch("api/collection.json"));

let collection = [];

async function $$main(request) {
	render("loading ...");
	collection = await prepareCollection(request);

	$$display();
}

function $$display() {
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
}

async function prepareCollection(promiseResponse) {
	const collection = await promiseResponse.then(rep => rep.json());
	prepareItem(collection);
	return collection;
}

function prepareItem(parent) {
	parent.children ??= [];
	parent.requests ??= [];
	for (const item of parent.children) {
		item[PARENT] = parent;
		prepareItem(item);
	}
	for (const request of parent.requests) {
		request[PARENT] = parent;
		defineHeadPropertie(request);
	}
}

function resetCollection() {
	collection = { name: "Collection" };
	prepareItem(collection);
	$$display();
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
			" ", $b("button.execute", "-", () => {
				if (!col[PARENT]) {
					resetCollection();
				} else {
					const children = col[PARENT].children;
					children.splice(children.indexOf(col), 1);
					$$display();
				}
			}),
		),
		(col.children || []).map($treeItem),
		(col.requests || []).map($treeRequest),
	);
}

function $treeRequest(request) {
	return {
		$: "div.tree-item.item-request",
		c: [
			$("span.tree-name", request.name),
			" ",
			$("span.status",
				!request.response ? "[?]" : !request.response.testFails ? "[V]" : "[X]"
			),
			" ", $b("button.execute", "-", () => {
				const parentRequests = request[PARENT].requests;
				parentRequests.splice(parentRequests.indexOf(request), 1);
				$$display();
			}),
		],
		onclick() { $$content(request); },
	};
}
