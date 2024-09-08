

const outputField = document.getElementById('output');
const inputField = document.getElementById('input');
const submitButton = document.getElementById('submit');

inputField.addEventListener('keydown', (event) => {
    if (event.key === 'Enter') {
        event.preventDefault();
        submit();
    }
});
submitButton.addEventListener('click', () => submit());

function submit() {
    console.log('Submitting');

    const input = inputField.value;
    inputField.value = '';
    outputField.innerText += '\n> ' + input + '\n';


    fetch('/input', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            content: input,
        })
    })
        .then(response => response.json())
        .then(data => {
            console.log(data);
            outputField.innerText += data.description;
        })
        .catch(error => {
            console.error(error);
        });
}
