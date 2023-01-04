const chokidar = require('chokidar');
const kill = require('tree-kill');
const spawn = require('child_process').spawn;

const [, , program, ...pargs] = process.argv

console.info(`:: Executing: ${program} ${pargs.join(' ')}`)
let child = undefined;
let processCoolDown = false
let stdinFiles = ''

function spawnChild() {
    processCoolDown = true
    child = spawn(program, pargs, {
        stdio: 'inherit'
    })
    setTimeout(() => {
        processCoolDown = false
    }, 1000)
}

spawnChild()

process.stdin.on('data', data => stdinFiles += data)
process.stdin.on('end', () => {
    filesToWatch = stdinFiles.split(/\s/).map(e => e.trim()).filter(e => e)
    filesToWatch.forEach(e => {
        console.log(":: Watching:", e)
    })
    let fileSet = new Set(filesToWatch)
    chokidar.watch(filesToWatch).on('change', (path) => {
        if (!fileSet.has(path)) return;
        if (processCoolDown) {
            console.log(`:: Changed ${path}. Waiting for old process to finish cooldown...`);
            return
        }
        console.log(`:: Changed ${path}. Killing old process...`);

        if (child) {
            kill(child.pid)
        }
        spawnChild()

    });
})
