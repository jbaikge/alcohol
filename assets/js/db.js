var Database = (function($) {
	var Database = {
		_db: {}
	}

	var idb = window.indexedDB || window.webkitIndexedDB || window.mozIndexedDB || window.moz_indexedDB
	var upgrading = false

	Database.Open = function(name, version, callback) {
		var request = idb.open(name, version)
		request.onupgradeneeded = onUpgradeNeeded
		request.onsuccess = function(e) {
			Database._db = e.target.result
			console.debug("Open Success %d", Database._db.version)
			if (!upgrading) {
				console.debug("Calling callback in success")
				callback(Database)
			}
		}
		request.onerror = onError
		return Database
	}

	var applyVersion = function(db, version) {
		console.debug("Downloading instructions for version %d", version)
		$.getJSON("/upgrade/" + version, function(data) {
			var transaction = db.transaction(["prices"], "readwrite")
			transaction.oncomplete = function(e) {
				console.debug("Calling callback in transaction complete")
				upgrading = false
			}
			transaction.onerror = onError

			var store = transaction.objectStore("prices")
			for (var i in data) {
				var request = store.add(data[i])
				request.onsuccess = function(e) {}
				request.onerror = onError
			}
		})
	}

	var initStores = function(db) {
		console.debug("Initializing stores")
		var prices = db.createObjectStore("prices", {
			keyPath: "Sku"
		})
		prices.createIndex("Type", "Type", {
			unique: false
		})
	}

	var onError = function(e) {
		console.debug("There was an error:", e.target)
	}

	var onUpgradeNeeded = function(e) {
		console.debug("Ugrade needed, old version: %d, new version: %d", e.oldVersion, e.newVersion)
		var db = e.target.result
		Database._upgrading = true
		e.target.transaction.onerror = onError

		if (e.oldVersion <= 1) {
			initStores(db)
		}

		applyVersion(db, e.newVersion)
	}

	return Database
}(jQuery))

function cb(db) {
	db._db.transaction(["prices"]).objectStore("prices").openCursor().onsuccess = function(e) {
		var cursor = e.target.result
		if (cursor) {
			console.debug("%s: %s", cursor.value.Sku, cursor.value.Name)
			cursor.continue()
		} else {
			console.debug("Reached the end")
		}
	}
}
