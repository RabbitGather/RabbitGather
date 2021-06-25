// Generated by CoffeeScript 2.5.1
(function() {
  var ref,
    boundMethodCheck = function(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new Error('Bound instance method accessed before binding'); } };

  ref = window.ServerFile = (function() {
    class ServerFile extends Backbone.Model {
      constructor() {
        super(...arguments);
        this.updateFileType = this.updateFileType.bind(this);
      }

      initialize() {
        this.on("change:type", this.updateFileType);
        this.updateFileType();
        if (this.get("dateCreated") === null) {
          return this.set("dateCreated", new Date());
        }
      }

      updateFileType() {
        boundMethodCheck(this, ref);
        if (this.get("type")) {
          return this.set("fileType", ServerFile.rawTypeToFileType(this.get("type")));
        } else {
          return this.set("fileType", ServerFile.filenameToFileType(this.get("name")));
        }
      }

      static rawTypeToFileType(rawType) {
        if (rawType.indexOf("image") !== -1) {
          return ServerFile.fileTypeEnum.IMG;
        }
        if (rawType.indexOf("html") !== -1 || rawType === "text/plain") {
          return ServerFile.fileTypeEnum.HTML;
        }
        if (rawType.indexOf("css") !== -1) {
          return ServerFile.fileTypeEnum.CSS;
        }
        if (rawType.indexOf("handlebars") !== -1) {
          return ServerFile.fileTypeEnum.TEMPLATE;
        }
        if (rawType.indexOf("javascript") !== -1) {
          return ServerFile.fileTypeEnum.JS;
        }
        return console.error("Unable to identify file type: " + rawType);
      }

      static filenameToFileType(filename) {
        var ext;
        ext = filename.replace(/.*\.([a-z]+$)/i, "$1");
        switch (ext) {
          case "html":
            return ServerFile.fileTypeEnum.HTML;
          case "jpg":
          case "jpeg":
          case "png":
            return ServerFile.fileTypeEnum.IMG;
          case "css":
            return ServerFile.fileTypeEnum.CSS;
          case "js":
            return ServerFile.fileTypeEnum.JS;
          case "handlebars":
            return ServerFile.fileTypeEnum.TEMPLATE;
          case "tmpl":
            return ServerFile.fileTypeEnum.TEMPLATE;
        }
        console.error("Attempt to convert file name " + filename + " to file type failed, returning null.");
        return null;
      }

    };

    ServerFile.prototype.defaults = {
      name: "",
      size: 0,
      contents: "",
      type: "",
      fileType: "",
      isProductionVersion: false,
      isRequired: false,
      dateCreated: null,
      hasBeenEdited: false
    };

    // If these are updated, note that there are dependencies on the client-side
    // for rendering.
    ServerFile.fileTypeEnum = {
      HTML: "HTML",
      CSS: "CSS",
      JS: "JS",
      IMG: "IMG",
      TEMPLATE: "TEMPLATE",
      NONE: "NONE"
    };

    ServerFile.fileTypeToFileExt = {
      HTML: "html",
      CSS: "css",
      JS: "js",
      TEMPLATE: "handlebars"
    };

    ServerFile.fileExtToFileType = _.invert(ServerFile.fileTypeToFileExt);

    return ServerFile;

  }).call(this);

}).call(this);
