document.getElementById('exportBtn').addEventListener('click', async () => {
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

    const results = await chrome.scripting.executeScript({
        target: { tabId: tab.id },
        func: () => {
            return document.body.innerHTML;
        }
    });

    console.log(`Job description HTML: ${results[0].result}`);
    const html = results[0].result;
    const url = tab.url;

    try {
        copyToClipboard(`HTML: ${html}\n'mf-URL: ${url}`);
        console.log('HTML copied to clipboard');
    } catch (error) {
        console.error('Failed to copy HTML to clipboard:', error);
    }
});

function copyToClipboard(text) {
    const textArea = document.createElement("textarea");
    textArea.value = text;
    textArea.style.position = "fixed";
    textArea.style.left = "-999999px";
    textArea.style.top = "-999999px";
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    
    try {
        document.execCommand('copy');
        console.log('Copied to clipboard successfully');
    } catch (err) {
        console.error('Failed to copy: ', err);
    }
    
    document.body.removeChild(textArea);
}
