export function $(name, ...components) {
	return { $: name, c: components };
}

export function render(components, root = document.body) {
	if (typeof root == "string") return render(components, document.getElementById(root));
	root.innerHTML = "";
	const ids = { $: [] };
	createChildren(root, [components], ids);
	ids.$.forEach(f => f());
	delete ids.$;
	return ids;
}

function createChildren(parent, components, ids) {
	for (const c of components) {
		if (!c || c === true) {
			continue;
		} else if (Array.isArray(c)) {
			createChildren(parent, c, ids);
			continue;
		} else if (typeof c == "function") {
			createChildren(parent, [c(parent)], ids);
		} else if (typeof c == "string") {
			parent.innerText += c;
		} else if (c instanceof HTMLElement) {
			parent.append(c);
		} else {
			createElement(parent, c, ids);
		}
	}
}

const SPLITER = /([.#])([^.#]*)/g;

function createElement(parent, c, ids) {
	const element = document.createElement(c.$.split(SPLITER)[0]);
	element.innerText = c.t || "";

	// Attribute
	for (const [attr, attrType, attrValue] of c.$.matchAll(SPLITER)) {
		switch (attrType) {
			case '.':
				element.classList.add(attrValue);
				break;
			case '#':
				element.id = attrValue;
				ids[attrValue] = element;
				break;
		}
	}

	Object.keys(c)
		.filter(key => key.length > 1)
		.forEach(key => element[key] = c[key]);

	// Children
	createChildren(element, c.c ?? [], ids);

	// Add in the dom
	parent.append(element);

	// Will call the function c.f after the render.
	if (c.f) ids.$.push(() => c.f(element));
}
