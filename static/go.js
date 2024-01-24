document.querySelector('#go').addEventListener('click', (e) => {
    const url = document.querySelector("#url").value;
    const match = url.match(/([0-9]+)$/);

	if (match === null) {
		const error = document.querySelector("#error");
		error.innerText = "Invalid UG URL";
		return;
	}

    const tabId = match[1];
    location = `/${tabId}`;
});
