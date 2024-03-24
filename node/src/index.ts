
console.log("Hello world!");


// avoid global variable
const submmitFunc = (function () {
    let submmited = false;

    return function () {
        if (submmited) {
            console.log("already submmited");
        }
        submmited = true;
    }
}());

submmitFunc();
submmitFunc();