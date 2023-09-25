const serverUrl = "http://127.0.0.1:8080/currencies.json";
const dataUpdateIntervalMilliseconds = 300000;
const decimalDigitsInResultValue = 4;

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
    } else {
        info.text("");
    }

    if (!data && !storedData) {
        return;
    }

    data.unshift(rubleCurrency);

    if (isFirstTimeGetUpdated) {
        initPageElements(data);

        isFirstTimeGetUpdated = false;
    }

    initClicksListener(data);

    storedData = data;
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

const initPageElements = (data) => {
    const firstCurrency = data[0];
    const secondCurrency = data[1];

    leftCurrencyButton.text(
        getExtendedCurrencyName(firstCurrency));
    rightCurrencyButton.text(
        getExtendedCurrencyName(secondCurrency));
    
    leftRatio = firstCurrency.ratio;
    rightRatio = secondCurrency.ratio;

    calculateResult();

    fillList(leftCurrencyList, leftCurrencySide, data);
    fillList(rightCurrencyList, rightCurrencySide, data);
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

        element =
            "<li><a class=\"dropdown-item\" onclick=\"side=" +
            side + ";index=" + i + ";\">" +
            getExtendedCurrencyName(currency) +
            "</a></li>";

        list.append(element);

        i++;
    }
}

const initClicksListener = (data) => {
    $(".dropdown-item").on("click", () => {
        const currency = data[index];
        const currencyNameWithCharCode =
            getExtendedCurrencyName(currency);

        if (side === leftCurrencySide) {
            leftRatio = parseFloat(currency.ratio);
            leftCurrencyButton.text(currencyNameWithCharCode);
        } else {
            rightRatio = parseFloat(currency.ratio);
            rightCurrencyButton.text(currencyNameWithCharCode);
        }

        calculateResult();
    });
}

const calculateResult = () => {
    const ratio = rightRatio / leftRatio;

    result.text(
        ratio.toFixed(decimalDigitsInResultValue).toString());
}

const getExtendedCurrencyName = (currency) => {
    return currency.name + " (" + currency.charCode + ")";
}

main();
