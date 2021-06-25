// Generated by CoffeeScript 2.5.1
(function() {
  window.ClientServer = class ClientServer {
    constructor(options) {
      this.channelOnReady = this.channelOnReady.bind(this);
      this.channelOnUnavailableID = this.channelOnUnavailableID.bind(this);
      this.channelOnInvalidID = this.channelOnInvalidID.bind(this);
      this.readDesiredServerIDFromURL = this.readDesiredServerIDFromURL.bind(this);
      this.channelOnConnection = this.channelOnConnection.bind(this);
      this.channelOnConnectionClose = this.channelOnConnectionClose.bind(this);
      this.channelConnectionOnData = this.channelConnectionOnData.bind(this);
      this.setUpReceiveEventCallbacks = this.setUpReceiveEventCallbacks.bind(this);
      this.sendEventTo = this.sendEventTo.bind(this);
      // The client-browser requested a file or path that resulted in an error.
      // The file might not exist, or evaluating the path may result in an error
      //  due to the user writing broken server code for the path.
      this.sendFailure = this.sendFailure.bind(this);
      this.serveFile = this.serveFile.bind(this);
      this.parsePath = this.parsePath.bind(this);
      // Returns the contents for the given path with the params.
      // foundRoute is an optional parameter that must be the corresponding dynamic path
      //   if the path is a dynamic path (ie, if the path is in the routecollection), or null
      //   if the path is a static file.
      // Returns either the html string, or null if none can be found.

      // TODO (?) handle leading slash and handle "./file"
      this.getContentsForPath = this.getContentsForPath.bind(this);
      this.recordResourceRequest = this.recordResourceRequest.bind(this);
      this.onResourceChange = this.onResourceChange.bind(this);
      this.removePeerFromClientBrowserResourceRequests = this.removePeerFromClientBrowserResourceRequests.bind(this);
      this.onDBChange = this.onDBChange.bind(this);
      this.serverFileCollection = options.serverFileCollection;
      this.routeCollection = options.routeCollection;
      this.appView = options.appView;
      this.userDatabase = options.userDatabase;
      this.desiredServerID = this.readDesiredServerIDFromURL();
      this.eventTransmitter = new EventTransmitter();
      this.userSessions = new UserSessions();
      this.dataChannel = new ClientServerDataChannel({
        onConnectionCallback: this.channelOnConnection,
        onDataCallback: this.channelConnectionOnData,
        onReady: this.channelOnReady,
        onConnectionCloseCallback: this.channelOnConnectionClose,
        desiredServerID: this.desiredServerID,
        onUnavailableIDCallback: this.channelOnUnavailableID,
        onInvalidIDCallback: this.channelOnInvalidID
      });
      this.setUpReceiveEventCallbacks();
      this.clientBrowserConnections = {};
      // TODO ensure that this is false
      this.isPushChangesEnabled = false;
      if (this.isPushChangesEnabled) {
        this.clientBrowserResourceRequests = {};
        this.serverFileCollection.on("change", this.onResourceChange);
        this.routeCollection.on("change", this.onResourceChange);
        this.userDatabase.on("onDBChange", this.onDBChange);
      }
    }

    channelOnReady() {
      var serverID;
      serverID = this.dataChannel.id;
      this.appView.trigger("setServerID", serverID);
      this.serverFileCollection.initLocalStorage(serverID);
      this.routeCollection.initLocalStorage(serverID);
      return this.userDatabase.initLocalStorage(serverID);
    }

    channelOnUnavailableID() {
      return this.appView.trigger("onUnavailableID", this.desiredServerID);
    }

    channelOnInvalidID() {
      return this.appView.trigger("onInvalidID", this.desiredServerID);
    }

    readDesiredServerIDFromURL() {
      if (/\/server\//.test(location.pathname)) {
        return location.pathname.replace(/\/server\//, "");
      }
      return null;
    }

    channelOnConnection(connection) {
      var contents, foundRoute, landingPage;
      landingPage = this.serverFileCollection.getLandingPage();
      // connection.peer is the socket id of the remote peer this connection
      // is connected to
      this.clientBrowserConnections[connection.peer] = connection;
      this.userSessions.addSession(connection.peer);
      this.appView.updateConnectionCount(_.size(this.clientBrowserConnections));
      foundRoute = this.routeCollection.findRouteForPath("/index");
      // Check if path mapping or a static file for /index exists -- otherwise send index.html
      if (foundRoute !== null && foundRoute !== void 0) {
        contents = this.getContentsForPath("/index", {}, foundRoute, connection.peer);
        if (contents && !contents.error) {
          landingPage = {
            fileContents: contents.result,
            filename: "index",
            type: "text/html"
          };
        }
      }
      return this.eventTransmitter.sendEvent(connection, "initialLoad", landingPage);
    }

    channelOnConnectionClose(connection) {
      if (connection && connection.peer) {
        this.userSessions.removeSession(connection.peer);
      }
      if (connection && connection.peer && _.has(this.clientBrowserConnections, connection.peer)) {
        delete this.clientBrowserConnections[connection.peer];
      }
      this.appView.updateConnectionCount(_.size(this.clientBrowserConnections));
      return this.removePeerFromClientBrowserResourceRequests(connection.peer);
    }

    channelConnectionOnData(data) {
      return this.eventTransmitter.receiveEvent(data);
    }

    setUpReceiveEventCallbacks() {
      return this.eventTransmitter.addEventCallback("requestFile", this.serveFile);
    }

    sendEventTo(socketId, eventName, data) {
      var connection;
      connection = this.clientBrowserConnections[socketId];
      return this.eventTransmitter.sendEvent(connection, eventName, data);
    }

    sendFailure(data, errorMessage) {
      var page404, response;
      if (data.type === "ajax") {
        response = {
          fileContents: "",
          type: data.type,
          textStatus: "error",
          errorThrown: errorMessage,
          requestId: data.requestId
        };
      } else {
        page404 = this.serverFileCollection.get404Page();
        response = {
          filename: page404.filename,
          fileContents: page404.fileContents,
          fileType: page404.type,
          type: data.type,
          errorMessage: errorMessage
        };
      }
      return this.sendEventTo(data.socketId, "receiveFile", response);
    }

    serveFile(data) {
      var contents, extraParams, fileType, foundRoute, foundServerFile, name, paramData, path, rawPath, response, slashedPath, val;
      rawPath = data.filename || "";
      if (_.isObject(rawPath)) { // In case the user passed an object with a url field instead.
        rawPath = rawPath.url;
      }
      [path, paramData] = this.parsePath(rawPath);
      if (data.options && data.options.data) { // Happens for ajax and form submits
        // Merge in any extra parameters passed with the ajax request.
        if (typeof data.options.data === "string") {
          extraParams = URI.parseQuery(paramData); // Return object mapping of get params in data.options.data
        } else {
          extraParams = data.options.data;
        }
        for (name in extraParams) {
          val = extraParams[name];
          paramData[name] = val;
        }
      }
      slashedPath = "/" + path;
      foundRoute = this.routeCollection.findRouteForPath(slashedPath);
      foundServerFile = this.serverFileCollection.findWhere({
        name: path,
        isProductionVersion: true
      });
      // Check if path mapping or a static file for this path exists -- otherwise send failure
      if ((foundRoute === null || foundRoute === void 0) && !this.serverFileCollection.hasProductionFile(path)) {
        console.error("Error: Client requested " + rawPath + " which does not exist on server.");
        this.sendFailure(data, "Not found");
        return;
      }
      if (foundRoute === null || foundRoute === void 0) {
        fileType = this.serverFileCollection.getFileType(path);
      } else {
        fileType = "UNKNOWN";
      }
      contents = this.getContentsForPath(path, paramData, foundRoute, data.socketId);
      // Check if following the path results in valid contents -- otherwise send failure
      if (!contents || contents.error) {
        console.error("Error: Function evaluation for  " + rawPath + " generated an error, returning 404: " + contents.error);
        this.sendFailure(data, "Internal server error");
        return;
      }
      // if contents.result and contents.result.extra is "redirect"  # Option to also return a function to be executed
      //   contents.result.fcn()
      // Construct the response to send with the contents
      response = {
        filename: rawPath,
        fileContents: contents.result,
        type: data.type,
        fileType: fileType
      };
      if (data.type === "ajax") {
        response.requestId = data.requestId;
      } else {
        this.recordResourceRequest(data.socketId, data, foundServerFile, foundRoute);
      }
      return this.sendEventTo(data.socketId, "receiveFile", response);
    }

    parsePath(fullPath) {
      var params, pathDetails;
      if (!fullPath || fullPath === "") {
        return ["", {}];
      }
      pathDetails = URI.parse(fullPath);
      params = URI.parseQuery(pathDetails.query);
      return [pathDetails.path, params];
    }

    getContentsForPath(path, paramData, foundRoute, socketId) {
      var match, runRoute, slashedPath;
      if (foundRoute === null || foundRoute === void 0) {
        return {
          "result": this.serverFileCollection.getContents(path)
        };
      }
      // Otherwise, handle a dynamic path
      slashedPath = "/" + path;
      // TODO (?) flesh out with params, etc.
      match = slashedPath.match(foundRoute.pathRegex);
      runRoute = foundRoute.getExecutableFunction(paramData, match.slice(1), this.serverFileCollection.getContents, this.userDatabase.database, this.userSessions.getSession(socketId));
      return runRoute();
    }

    recordResourceRequest(peerID, data, foundServerFile, foundRoute) {
      var resource, resourceName;
      if (!this.isPushChangesEnabled) {
        return;
      }
      if (!((foundServerFile && foundServerFile.get("fileType") === ServerFile.fileTypeEnum.HTML) || foundRoute)) {
        return;
      }
      resource = null;
      if (foundServerFile) {
        resource = foundServerFile;
      } else if (foundRoute) {
        resource = foundRoute;
      }
      if (!resource) {
        return;
      }
      resourceName = resource.get("name");
      this.removePeerFromClientBrowserResourceRequests(peerID);
      if (!this.clientBrowserResourceRequests[resourceName]) {
        this.clientBrowserResourceRequests[resourceName] = [];
      }
      return this.clientBrowserResourceRequests[resourceName].push({
        peerID: peerID,
        data: data
      });
    }

    onResourceChange(resource) {
      var interestedPeers;
      if (!this.isPushChangesEnabled) {
        return;
      }
      if (!resource.get("isProductionVersion")) {
        return;
      }
      interestedPeers = this.clientBrowserResourceRequests[resource.get("name")];
      return _.each(interestedPeers, (interestedPeer) => {
        return this.serveFile(interestedPeer.data);
      });
    }

    removePeerFromClientBrowserResourceRequests(peerID) {
      var resourceNames;
      if (!this.isPushChangesEnabled) {
        return;
      }
      resourceNames = _.keys(this.clientBrowserResourceRequests);
      return _.each(resourceNames, (resourceName) => {
        return this.clientBrowserResourceRequests[resourceName] = _.filter(this.clientBrowserResourceRequests[resourceName], (interestedPeer) => {
          return interestedPeer.peerID !== peerID;
        });
      });
    }

    onDBChange() {
      var resourceNames;
      if (!this.isPushChangesEnabled) {
        return;
      }
      resourceNames = _.keys(this.clientBrowserResourceRequests);
      return _.each(resourceNames, (resourceName) => {
        var interestedPeers, route;
        // This is a bit of a hack -- not currently used experiment for push-on-DB-change functionality.
        route = this.routeCollection.findWhere({
          name: resourceName,
          isProductionVersion: true
        });
        if (route && /database\.insert\(|database\(.*?\)\.remove\(|database\(.*?\)\.update\(/.test(route.get("routeCode"))) {
          return;
        }
        interestedPeers = this.clientBrowserResourceRequests[resourceName];
        return _.each(interestedPeers, (interestedPeer) => {
          return this.serveFile(interestedPeer.data);
        });
      });
    }

  };

}).call(this);
