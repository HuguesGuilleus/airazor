import { render, $, $b } from "./util.mjs";
import { $$content, defineHeadPropertie } from "./content.mjs";

const PARENT = Symbol("parent");

let collection = [];
let focusRequest = null;

$$main(fetch("api/collection.json"));

async function $$main(request) {
	render("loading ...");
	collection = await prepareCollection(request);

	$$display();
}

function $$display() {
	render([
		$("aside.tree",
			$b("button.bigAct", "Save", "", () => {
				fetch("/api/collection.json", {
					method: "PUT",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify(collection)
				}).then(({ ok }) => alert("saving: " + (ok ? "SUCESS" : "FAIL")))
			}),
			$b("button.bigAct", "Run all", "", () => {
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
	if (focusRequest) {
		$$content(focusRequest);
	}
}

async function prepareCollection(promiseResponse) {
	focusRequest = null;
	const collection = await promiseResponse.then(rep => rep.json());
	prepareItem(collection);
	return collection;
}

function prepareItem(parent) {
	parent.children ??= [];
	parent.requests ??= [];
	sortByName(parent.children);
	sortByName(parent.requests);
	for (const item of parent.children) {
		item[PARENT] = parent;
		prepareItem(item);
	}
	for (const request of parent.requests) {
		request[PARENT] = parent;
		defineHeadPropertie(request);
	}
	if (!focusRequest) {
		focusRequest = parent.requests[0];
	}
}
function sortByName(items) {
	items.sort((a, b) => a.name.toLowerCase() > b.name.toLowerCase());
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
			" ", $b("button.execute", "C+", "New Collection", () => {
				col.children.push({ name: prompt("Name of the new collection") });
				prepareItem(col);
				$$display();
			}),
			" ", $b("button.execute", "R+", "New Request", () => {
				col.requests.push({ name: prompt("Name of the new request") });
				prepareItem(col);
				$$display();
			}),
			$renameItem(col, false),
			" ", $b("button.execute", "-", "Remove the colection", () => {
				if (!col[PARENT]) {
					resetCollection();
				} else {
					const children = col[PARENT].children;
					children.splice(children.indexOf(col), 1);
					$$display();
				}
			}),
		),
		col.children.map($treeItem),
		col.requests.map($treeRequest),
	);
}

function $treeRequest(request) {
	return {
		$: ["div.tree-item.item-request", request == focusRequest && ".item-focus"],
		c: [
			$("span.tree-name", request.name),
			" ", $("span.status",
				!request.response ? "[?]" : !request.response.testFails ? "[V]" : "[X]"
			),
			$renameItem(request, true),
			" ", $b("button.execute", "-", "Remove the request", () => {
				const parentRequests = request[PARENT].requests;
				parentRequests.splice(parentRequests.indexOf(request), 1);
				$$display();
			}),
		],
		onclick() {
			focusRequest = request;
			$$display();
			$$content(request);
		},
	};
}

function $renameItem(item, isRequest) {
	return [
		" ",
		$b("button.execute", "R", "Rename", () => {
			const name = window.prompt("Rename", item.name);
			if (!name) return;
			item.name = name;
			if (isRequest) {
				focusRequest = item;
			}
			prepareItem(col);
			$$display();
		}),
	];
}
