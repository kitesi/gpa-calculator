// script file to generate grade files for testing

function getTemplate() {
    const midTermMax = randomInt(40, 50);
    const finalMax = randomInt(40, 50);

    const template = `~ Meta
> Homework
weight = 0.2
data = 20/20,  ${random(12, 20)}/20, ${random(18, 20)}/20, ${random(10, 20)}/20

> Quizes
weight = 0.1
data = 20/20,  ${random(12, 20)}/20, ${random(18, 20)}/20, ${random(10, 20)}/20

> Mid Term 
weight = 0.3
data = ${random(20, midTermMax)}/${randomInt(40, midTermMax)}

> Final Eaxm
weight = 0.4
data =  ${random(30, finalMax)}/${finalMax}
`;
    return template;
}

function random(min, max) {
    return (Math.random() * (max - min + 1) + min).toFixed(2);
}

function randomInt(min, max) {
    return Math.floor(random(min, max));
}

function main() {
    if (!process.argv[2]) {
        console.error('No level provided');
        process.exit(1);
    }

    const fs = require('fs');
    const courses = ['ma', 'gov', 'cs', 'lang'];

    for (const course of courses) {
        fs.writeFile(
            course + process.argv[2] + '.grade',
            getTemplate(),
            (err) => {
                if (err) {
                    throw err;
                }
            }
        );
    }
}

main();
