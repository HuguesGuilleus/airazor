export default {
	name: "Global",
	children: [
		{
			name: "Sub1", requests: [
				{
					name: "Not found",
					method: "GET",
					URL: "https://jsonplaceholder.typicode.com/notfound",
					test: `equal(200, response.code);\nequal("{}", response.text);`,
					response: {}
				},
			]
		},
		{
			name: "Sub1", requests: []
		},
	],
};