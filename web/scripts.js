const serverUrl = "http://127.0.0.1:8080/currencies.json";
const dataUpdateIntervalMilliseconds = 300000;

const leftCurrencyButton = $("#left-currency-button");
const leftCurrencyList = $("#left-currency-list");
const rightCurrencyButton = $("#right-currency-button");
const rightCurrencyList = $("#right-currency-list");
const result = $("#result");
const info = $("#info-container");
const leftCurrencySide = 1;
const rightCurrencySide = 2;

const rubleCurrency = {
    name: "Российский рубль",
    charCode: "RUB",
    ratio: "1.0"
}

let storedData;
let index;
let side;
let leftRatio = 1;
let rightRatio = 1;
let isFirstTimeGetUpdated = true;

const main = () => {
    updateDataFromServerWithInterval(
        dataUpdateIntervalMilliseconds
    );
}

const updateDataFromServerWithInterval = (timeInterval) => {
    updateDataFromServer();

    return(setInterval(updateDataFromServer, timeInterval));
}

const updateDataFromServer = () => {
    const data = getDataFromServer(serverUrl);

    if (!data) {
        info.text("server is not responding");
    }

    info.text("");

    data.unshift(rubleCurrency);

    storedData = data;

    if (!storedData) {
        return;
    }

    if (isFirstTimeGetUpdated) {
        initPageElements(data);

        isFirstTimeGetUpdated = false;
    }

    initClicksListener();
}

const initPageElements = (data) => {
    const firstElement = data[0];
    const secondElement = data[1];

    leftCurrencyButton.text(firstElement.name);
    rightCurrencyButton.text(secondElement.name);
    
    leftRatio = firstElement.ratio;
    rightRatio = secondElement.ratio;

    calculateResult();

    fillList(leftCurrencyList, leftCurrencySide, data);
    fillList(rightCurrencyList, rightCurrencySide, data);
}

const initClicksListener = () => {
    $(".dropdown-item").on("click", () => {
        const currency = storedData[index];

        if (side === leftCurrencySide) {
            leftRatio = parseFloat(currency.ratio);
            leftCurrencyButton.text(currency.name);
        } else {
            rightRatio = parseFloat(currency.ratio);
            rightCurrencyButton.text(currency.name);
        }

        calculateResult();
    });
}

const getDataFromServer = (url) => {
    let result;

    console.debug("connecting to server...");

    $.ajax({
        url: url,
        type: "get",
        dataType: "json",
        async: false,
        cache: false,
        success: (data) => {
            result = data;
        }
    });

    return result;
}

const fillList = (list, side, data) => {
    const divider = "<li><hr class=\"dropdown-divider\"></li>";

    let currency;
    let element;
    let i = 0;

    while (i < data.length) {
        if (i === 1) {
            list.append(divider);    
        }

        currency = data[i];

        element = "<li><a class=\"dropdown-item\" href=\"#\" onclick=\"side=" +
            side + ";index=" + i + ";\">"  + currency.name + "</a></li>";

        list.append(element);

        i++;
    }
}

const calculateResult = () => {
    const ratio = rightRatio / leftRatio;

    result.text(ratio.toFixed(4).toString());
}

main();
