const btn = document.getElementById('btn');
const outputSSML = document.querySelector('#output-ssml');
btn.addEventListener('click', () => {
    console.log('Button clicked');
    const text = document.querySelector('#dialogue').value;
    const voiceA = document.querySelector('#voiceA').value;
    const voiceB = document.querySelector('#voiceB').value;

    const sentences = text.split('\n').map(sentence => sentence.trim()).filter(sentence => sentence);
    //eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGU1ZDE0YmEtMDc4Ny00NWNiLTlmOTEtZDlmZTkyNWU3ZDc0IiwidHlwZSI6ImFwaV90b2tlbiJ9.kq_iMKBsu-03o1QQEgj8tvbNn-Gk5etT-5sCl-ZuQyA
    // const options = {
    //     method: 'POST',
    //     url: 'http://api.edenai.run/v2/translation/language_detection',
    //     headers: {
    //         authorization: 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGU1ZDE0YmEtMDc4Ny00NWNiLTlmOTEtZDlmZTkyNWU3ZDc0IiwidHlwZSI6ImFwaV90b2tlbiJ9.kq_iMKBsu-03o1QQEgj8tvbNn-Gk5etT-5sCl-ZuQyA',
    //     },
    //     data: {
    //         providers: "amazon.google",
    //         text: sentences[0],
    //     }
    // };
    const lngDetector = new (require('languagedetect'));
    let detectedLang = lngDetector.detect(sentences[0], 1);
    console.log(detectedLang[0][0]);
    if (detectedLang[0][0] == "en") {
        detectedLang[0][0] = "en-US";
    } else {
        detectedLang[0][0] = "vi-VI";
    }
    let ssml = `<speak xml:lang=${detectedLang[0][0]}>`;
    sentences.forEach((sentence, index) => {
        const voice = index % 2 === 0 ? voiceA : voiceB;
        ssml += `
        <voice name="${voice}">
            ${sentence}
        </voice>`;
    });
    ssml += '</speak>';
    outputSSML.textContent = ssml;
});