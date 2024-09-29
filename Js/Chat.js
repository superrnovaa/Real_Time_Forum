import * as socket from "./WS.js";

export function ChatTemp() {
  let template = `
     <div class="PartnerInfo">
     <div class="PersonalPhoto" id="partnerProfileImg"><img src="../ProfileImages/${sessionStorage.getItem("partner-Img")}" /></div>
     <p class="Username" id="partnerUsername">${sessionStorage.getItem("partner-Username")}</p>  
        </div>
    <div class="ChatBody"></div>
      <div class="Message">
        <input id="ChatMessage" title="Write Message" tabindex="i" pattern="\d+"
          placeholder="Message.." class="MsgInput" type="text" />
        <div class="SendChatMessage">
          <svg xmlns="http://www.w3.org/2000/svg" version="1.0" width="30.000000pt" height="30.000000pt"
            viewBox="0 0 30.000000 30.000000" preserveAspectRatio="xMidYMid meet" class="SendSVG">
            <g transform="translate(0.000000,30.000000) scale(0.100000,-0.100000)" fill="#2b2b35" stroke="none">
              <path
                d="M44 256 c-3 -8 -4 -29 -2 -48 3 -31 5 -33 56 -42 28 -5 52 -13 52 -16 0 -3 -24 -11 -52 -16 -52 -9 -53 -9 -56 -48 -2 -21 1 -43 6 -48 10 -10 232 97 232 112 0 7 -211 120 -224 120 -4 0 -9 -6 -12 -14z">
              </path>
            </g>
          </svg>
        </div>
      </div>`;


  return template;
}


export async function renderChatTemplate() {
  //  const container =  `<div class="SinglePostContainer"></div>`
  const style = `<link href="../Style/Chat.css" rel="stylesheet" />`;
  // document.querySelector(".page").innerHTML = "";
  for (let i = 0; i < 10; i++) {
    if (document.querySelector(".page")) {
      document.querySelector(".page").innerHTML = style + ChatTemp();
      break;
    }
    await new Promise(resolve => setTimeout(resolve, 100));
  }

  // document.querySelector(".SinglePostContainer").innerHTML = ChatTemp();
  socket.renderChatEvents();
}

