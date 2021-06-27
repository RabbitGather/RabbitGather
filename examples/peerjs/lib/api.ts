import { util } from "./util";
import logger from "./logger";

export class API {
  constructor(private readonly _options: any) {
    // 把_options直接加進來作為變數
  }

  private _buildUrl(method: string): string {
    // https://0.peerjs.com/peerjs/id?ts=16246764265920.288844543657901
    const protocol = this._options.secure ? "https://" : "http://";
    let url =
      // https://
      protocol +
      // 0.peerjs.com
      this._options.host +
      ":" +
      // 443
      this._options.port +
      // /
      this._options.path +
      // peerjs
      this._options.key +
      "/" +
      // id
      method;
    const queryString = "?ts=" + new Date().getTime() + "" + Math.random();
    url += queryString;

    return url;
  }

  /** Get a unique ID from the server via XHR (XMLHttpRequest) and initialize with it. */
  async retrieveId(): Promise<string> {
    const url = this._buildUrl("id");

    try {
      const response = await fetch(url);

      if (response.status !== 200) {
        throw new Error(`Error. Status:${response.status}`);
      }

      return response.text();
    } catch (error) {
      logger.error("Error retrieving ID", error);

      let pathError = "";

      if (
        this._options.path === "/" &&
        this._options.host !== util.CLOUD_HOST
      ) {
        pathError =
          " If you passed in a `path` to your self-hosted PeerServer, " +
          "you'll also need to pass in that same path when creating a new " +
          "Peer.";
      }

      throw new Error("Could not get an ID from the server." + pathError);
    }
  }

  /** @deprecated */
  async listAllPeers(): Promise<any[]> {
    const url = this._buildUrl("peers");

    try {
      const response = await fetch(url);

      if (response.status !== 200) {
        if (response.status === 401) {
          let helpfulError = "";

          if (this._options.host === util.CLOUD_HOST) {
            helpfulError =
              "It looks like you're using the cloud server. You can email " +
              "team@peerjs.com to enable peer listing for your API key.";
          } else {
            helpfulError =
              "You need to enable `allow_discovery` on your self-hosted " +
              "PeerServer to use this feature.";
          }

          throw new Error(
            "It doesn't look like you have permission to list peers IDs. " +
              helpfulError
          );
        }

        throw new Error(`Error. Status:${response.status}`);
      }

      return response.json();
    } catch (error) {
      logger.error("Error retrieving list peers", error);

      throw new Error("Could not get list peers from the server." + error);
    }
  }
}
