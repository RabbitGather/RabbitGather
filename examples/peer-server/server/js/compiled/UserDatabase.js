// Generated by CoffeeScript 2.5.1
(function() {
  `Simple wrapper for a taffy database. 

Might be extended to back up the database to local storage, a zip file, etc.`;
  var ref,
    boundMethodCheck = function(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new Error('Bound instance method accessed before binding'); } };

  ref = window.UserDatabase = class UserDatabase extends Backbone.Model {
    constructor() {
      super(...arguments);
      this.initLocalStorage = this.initLocalStorage.bind(this);
      this.toString = this.toString.bind(this);
      this.fromJSONArray = this.fromJSONArray.bind(this);
      this.runQuery = this.runQuery.bind(this);
      this.clear = this.clear.bind(this);
      this.onDBChange = this.onDBChange.bind(this);
    }

    initialize() {
      this.database = TAFFY();
      return this.database.settings({
        onDBChange: this.onDBChange
      });
    }

    initLocalStorage(namespace) {
      boundMethodCheck(this, ref);
      this.database.store(namespace + "-UserDatabase");
      return this.trigger("initLocalStorage");
    }

    toString(pretty) {
      boundMethodCheck(this, ref);
      if (pretty) {
        return JSON.stringify(this.database().get(), null, 4);
      }
      return this.database().stringify();
    }

    fromJSONArray(array) {
      boundMethodCheck(this, ref);
      return this.database.insert(array);
    }

    runQuery(query) {
      var code;
      boundMethodCheck(this, ref);
      code = "(function(database) { " + query + " }).call(null, this.database)";
      return eval(code);
    }

    clear() {
      boundMethodCheck(this, ref);
      return this.database().remove();
    }

    onDBChange() {
      boundMethodCheck(this, ref);
      return this.trigger("onDBChange");
    }

  };

}).call(this);
