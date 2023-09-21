const { app, BrowserWindow, ipcMain } = require('electron')
const path = require('node:path');
const { spawn } = require('node:child_process');

const createWindow = () => {
    const win = new BrowserWindow({
        width: 800,
        height: 600,
    })

    win.loadURL('http://localhost:6969')
}

app.whenReady().then(() => {
    ipcMain.handle('ping', () => 'pong')
    createWindow()

    const goNodeArgs = [
        "-port",
        "3000",
        "-nick",
        "ajmal",
        "-room",
        "ajmal",
    ]

    let goNode = spawn(
        path.join(app.getAppPath(), "src", "node", "node.exe"),
        goNodeArgs,
        {
            detached: false,
        },
    )
    goNode.stderr.on("data", (d) =>
        console.log("chat-gonode:", d.toString()),
    )

    app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) {
            createWindow()
        }
    })
})

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit()
    }
})