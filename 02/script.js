const btn = document.getElementById('btn');
const copyBtn = document.getElementById('copy-btn');
const outputSSML = document.querySelector('#output-ssml');
btn.addEventListener('click', () => {
    const text = document.querySelector('#dialogue').value;
    const voiceA = document.querySelector('#voiceA').value;
    const voiceB = document.querySelector('#voiceB').value;

    const sentences = text.split('\n').map(sentence => sentence.trim()).filter(sentence => sentence);
    console.log(sentences);
    const lngDetector = new (require('languagedetect'));
    let detectedLang = lngDetector.detect(sentences[0], 1);
    console.log(detectedLang[0][0]);
    if (detectedLang[0][0] == "en") {
        detectedLang[0][0] = "en-US";
    } else {
        detectedLang[0][0] = "vi-VI";
    }
    let ssml = `<speak version=\"1.0\" xmlns=\"http://www.w3.org/2001/10/synthesis\" xml:lang=\"" + detectedLang[0][0] + "\">\n`;
    let curVoiceID = 0;
    const voices = [voiceA, voiceB];
    let curDialogue = "";
    sentences.forEach((sentence) => {
        const components = sentence.split(':');
        console.log(components);
        if (components.length === 1) {
            curDialogue += sentence + ' ';
        } else {
            if (curDialogue !== "") {
                ssml += `\t<voice name="${voices[curVoiceID]}">${curDialogue.substring(0, curDialogue.length - 1)}</voice>\n`;
            curVoiceID = 1 - curVoiceID;
            }
            curDialogue = "";
            curDialogue += components[1] + ' ';
        }
        
    });
    if (curDialogue !== "") {
        ssml += `\t<voice name="${voices[curVoiceID]}">${curDialogue}</voice>\n`;
    }
    ssml += `</speak>`;
    outputSSML.textContent = ssml;
    console.log(ssml);
    document.querySelector('#ssml-display').className = 'readonly';

    const contentLength = ssml.length;
    const minWidth = 300; // Minimum width in pixels
    const maxWidth = 800; // Maximum width in pixels
    const width = Math.min(maxWidth, Math.max(minWidth, contentLength * 8)); // Adjust the multiplier as needed
    outputSSML.style.width = `${width}px`;
});

copyBtn.addEventListener('click', () => {
    const ssmlContent = outputSSML.textContent;
    navigator.clipboard.writeText(ssmlContent).then(() => {
        alert('SSML content copied to clipboard!');
    }).catch(err => {
        console.error('Failed to copy text: ', err);
    });
});