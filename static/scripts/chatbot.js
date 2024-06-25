const chatInput = document.querySelector(".chat-input textarea")
const sendChatBtn = document.querySelector(".chat-input span")
const chatbox = document.querySelector(".chatbox")

const url = window.location.href;
const params = new URLSearchParams(url.split('?')[1]);
const product = params.get('product');
//console.log(product)

let userMessage;
let history;


const createChatLi = (message, className) => {
    const chatLi = document.createElement('li')
    chatLi.classList.add("chat",className);
    let chatContent = className === "outgoing" ? `<p>${message}</p>` : `<span class="material-symbols-outlined">smart_toy</span><p>${message}</p>`
    chatLi.innerHTML = chatContent;
    return chatLi;
}

const generateResponse = async (incomingChatLi) => {
    if (!incomingChatLi) {
        console.error("Incoming chat element is null or undefined.");
        return;
    }
    const API_URL = "http://127.0.0.1:8000";
    // const API_URL = window.location.href;
    const messageElement = incomingChatLi.querySelector("p");
    if (!messageElement) {
        console.error("Message element not found within incoming chat element.");
        return;
    }
    

    let requestBody = {
        message: userMessage,
        history: history,
        product: product,
    };

    const requestOptions = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
    };


    try {
        console.log(requestBody)
        const res = await fetch(API_URL, requestOptions);
        if (!res.ok) {
            throw new Error("Network response was not ok! :C");
        }
        const data = await res.json();
        console.log(data);

        messageElement.textContent = data.message;
        history = data.history;
        // console.log(messageElement.textContent);
    } catch (error) {
        console.error('Error in catch:', error);
        messageElement.textContent = "Oops! Something went wrong.\nI can not assist you at the moment, consider contacting the human support team! :-)";
    }
}


const handleChat = () => {
    userMessage = chatInput.value.trim();
    if(!userMessage) return;

    chatbox.appendChild(createChatLi(userMessage, 'outgoing'));
    chatInput.value = "";

    setTimeout(async () => {
        const incomingChatLi = createChatLi("Thinking...", "incoming")
        // const ThinkingGif = document.createElement('div');
        // ThinkingGif.innerHTML = '<div class="tenor-gif-embed" data-postid="14292188" data-share-method="host" data-aspect-ratio="1.86" data-width="100%"><a href="https://tenor.com/view/typing-dots-waiting-loading-gif-14292188">Typing Dots GIF</a>from <a href="https://tenor.com/search/typing-gifs">Typing GIFs</a></div> <script type="text/javascript" async src="https://tenor.com/embed.js"></script>';
        chatbox.appendChild(incomingChatLi)
        await generateResponse(incomingChatLi);
    }, 600);
}

sendChatBtn.addEventListener("click", handleChat)
