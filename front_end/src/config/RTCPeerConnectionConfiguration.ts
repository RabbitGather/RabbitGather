// import Peer, {  RTCConfiguration } from "@/lib/peerjs";
const RTCPeerConnectionConfiguration = {
  iceServers: [
    {
      urls: [
        "stun:stun.l.google.com:19302",
        "stun:stun.l.google.com:19302",
        "stun.l.google.com:19302",
        "stun1.l.google.com:19302",
        "stun2.l.google.com:19302",
        "stun3.l.google.com:19302",
        "stun4.l.google.com:19302",
        "stun01.sipphone.com",
        "stun.ekiga.net",
        "stun.fwdnet.net",
        "stun.ideasip.com",
        "stun.iptel.org",
        "stun.rixtelecom.se",
        "stun.schlund.de",
        "stunserver.org",
        "stun.softjoys.com",
        "stun.voiparound.com",
        "stun.voipbuster.com",
        "stun.voipstunt.com",
        "stun.voxgratia.org",
        "stun.xten.com",

        "turn:turn01.hubl.in?transport=udp",
        "turn:turn02.hubl.in?transport=tcp",
      ],
    },
    {
      urls: "turn:0.peerjs.com:3478",
      username: "peerjs",
      credential: "peerjsp",
    },
    {
      urls: "turn:numb.viagenie.ca",
      username: "webrtc@live.com",
      credential: "muazkh",
    },
    {
      urls: "turn:192.158.29.39:3478?transport=udp",
      credential: "JZEOEt2V3Qb0y27GRntt2u2PAYA=",
      username: "28224511:1379330808",
    },
    {
      urls: "turn:192.158.29.39:3478?transport=tcp",
      credential: "JZEOEt2V3Qb0y27GRntt2u2PAYA=",
      username: "28224511:1379330808",
    },
    {
      urls: "turn:turn.bistri.com:80",
      credential: "homeo",
      username: "homeo",
    },
    {
      urls: "turn:turn.anyfirewall.com:443?transport=tcp",
      credential: "webrtc",
      username: "webrtc",
    },
  ],
  sdpSemantics: "unified-plan",
};

export default RTCPeerConnectionConfiguration;
