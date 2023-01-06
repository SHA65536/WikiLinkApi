const leftform = document.getElementById('left-search-form');
const rightform = document.getElementById('right-search-form');

const leftresults = document.getElementById('left-search-results');
const rightresults = document.getElementById('right-search-results');

const wikiendpointSearch = `/search?q=`
const wikiendpointRandom = `/random`

let selectedLeft = false
let selectedRight = false

leftform.addEventListener('submit', handleSubmitLeft);
rightform.addEventListener('submit', handleSubmitRight);

// Handles search of start article
async function handleSubmitLeft(event) {
    event.preventDefault();
    leftresults.innerHTML = '';
    selectedLeft = false;
    const inputVal = document.getElementById('left-search-input').value.trim();
    try {
        const results = await searchWikipedia(inputVal);
        displayResults(results, leftresults, "left-item")
    } catch (err) {
        console.log(err);
        alert('Failed to search wikipedia');
    }
}

// Handles search of end article
async function handleSubmitRight(event) {
    event.preventDefault();
    rightresults.innerHTML = '';
    selectedRight = false;
    const inputVal = document.getElementById('right-search-input').value.trim();
    try {
        const results = await searchWikipedia(inputVal);
        displayResults(results, rightresults, "right-item")
    } catch (err) {
        console.log(err);
        alert('Failed to search wikipedia');
    }
}

// Handles selection of search results
async function handleSelect(event, side) {
    if (side === "left-item") {
        if (selectedLeft != false) {
            selectedLeft.classList.toggle('selected-button');
        }
        selectedLeft = event.originalTarget;
        selectedLeft.classList.toggle('selected-button');
    } else {
        if (selectedRight != false) {
            selectedRight.classList.toggle('selected-button');
        }
        selectedRight = event.originalTarget;
        selectedRight.classList.toggle('selected-button');
    }
    console.log(side, event);
}


// Searches wikipedia for searchQuery and returns the results
async function searchWikipedia(searchQuery) {
    const response = await fetch(wikiendpointSearch + searchQuery);
    if (!response.ok) {
        throw Error(response.statusText);
    }
    const json = await response.json();
    return json;
}

// Displays the results on the given side
function displayResults(results, side, sideclass) {
    results.query.search.forEach(result => {
        const url = `https://en.wikipedia.org/?curid=${result.pageid}`;
        side.insertAdjacentHTML(
            "beforeend",
            `<button class="result-item ${sideclass}" onClick="handleSelect(event, '${sideclass}')">
                <h5 class="result-title">
                    <span>${result.title}</span><br>
                </h5>
                <span class="result-snippet">${result.snippet}</span><br>
            </button><br>`
        );
    });
}

// Generates random articles on both sides
async function RandomArticle() {
    leftresults.innerHTML = '';
    selectedLeft = false;
    rightresults.innerHTML = '';
    selectedRight = false;
    for (i = 0; i < 5; i++) {
        const response = await fetch(wikiendpointRandom);
        if (!response.ok) {
            throw Error(response.statusText);
        }
        const json = await response.json();
        const result = Object.values(json.query.pages)[0]
        const url = `https://en.wikipedia.org/?curid=${result.pageid}`;
        rightresults.insertAdjacentHTML(
            "beforeend",
            `<button class="result-item right-item" onClick="handleSelect(event, 'right-item')">
                <h5 class="result-title">
                    <span>${result.title}</span><br>
                </h5>
                <span class="result-snippet">${result.extract}</span><br>
            </button><br>`
        );
    }
    for (i = 0; i < 5; i++) {
        const response = await fetch(wikiendpointRandom);
        if (!response.ok) {
            throw Error(response.statusText);
        }
        const json = await response.json();
        const result = Object.values(json.query.pages)[0]
        const url = `https://en.wikipedia.org/?curid=${result.pageid}`;
        leftresults.insertAdjacentHTML(
            "beforeend",
            `<button class="result-item left-item" onClick="handleSelect(event, 'left-item')">
                <h3 class="result-title">
                    <span>${result.title}</span><br>
                </h3>
                <span class="result-snippet">${result.extract}</span><br>
            </button><br>`
        );
    }
}

// Goes to the path search page with current selected stuff
function Search() {
    if (selectedLeft === false || selectedRight === false) {
        return
    }
    const dst = selectedLeft.childNodes[1].childNodes[1].textContent
    const src = selectedRight.childNodes[1].childNodes[1].textContent
    window.location.href = `/result?start=${src}&end=${dst}`
}