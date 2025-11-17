function sayHi(user) {
    console.log(`Js module say Hello, ${user}!`);
}

function test() {
    return "passed";
}

module.exports = {
    test: test,
    sayHi: sayHi
}
