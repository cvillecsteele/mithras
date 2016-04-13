// MITHRAS: Javascript configuration management tool for AWS.
// Copyright (C) 2016, Colin Steele
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//                  (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//              GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
(function(window) {
    'use strict'

    var sprintf = require("sprintf").sprintf;

    var Cache = function Cache(path) {
	this.path = path || "cache";
    };

    Cache.prototype = {
	init: function () {

	    // TODO: clean up old entries

	    var error = fs.mkdirAll(this.path, 0777);
	    if (error) {
		log(sprintf("Error creating cache dir '%s'", this.path));
		os.exit(1);
	    }

	    this.metaPath = filepath.join(this.path, ".meta");
	    var stat = fs.stat(this.metaPath);
	    if (stat && stat.Err == 2) {
		var result = fs.create(this.metaPath);
		if (result[1]) {
		    log(sprintf("Error creating cache: %s", result[1]));
		    os.exit(1);
		}
		fs.close(result[0]);
	    }

	    var result = fs.read(this.metaPath);
	    var error = result[1];
	    if (error) {
		log("Error reading cache metadata");
		os.exit(1);
	    }
	    this.exMap = {};
	    if (result[0]) {
		this.exMap = JSON.parse(result[0]);
	    }

	    return this;
	}

	get: function (key) {
	    // Is it cached?
	    if (!this.exMap[key]) {
		return null;
	    }
	    var path = filepath.join(this.path, key, "value");
	    var stat = fs.stat(path);
	    var mod = Date.parse(stat.ModTime);

	    // Handle expiration
	    var expiresAt = mod + (this.exMap[key] * 1000);
	    if (Date.now() > expiresAt) {
		delete this.exMap[key];
		this.writeMeta();
		var path = filepath.join(this.path, key);
		fs.removeAll(path);
		return null;
	    }

	    var result = fs.read(path);
	    var error = result[1];
	    if (error) {
		if (error.Err == 2) { // ENOENT
		    return null;
		}
		log(sprintf("Error reading cached file: %s", 
			    JSON.stringify(error)));
		os.exit(1);
	    }
	    return JSON.parse(result[0]);
	}

	writeMeta: function() {
	    error = fs.write(this.metaPath, JSON.stringify(this.exMap, null, 2), 0644);
	    if (error) {
		log(sprintf("Cache metadata write file error: %s", error));
		os.exit(1);
	    }
	}

	put: function(key, value, expiry) {
	    var path = filepath.join(this.path, key);

	    // Create dir & file
	    var error = fs.mkdirAll(path, 0777)
	    if (error) {
		log("Cache create dir error: %s", error);
		os.exit(1);
	    }
	    path = filepath.join(path, "value");
	    var result = fs.create(path)
	    var error = result[1];
	    if (error) {
		log(sprintf("Cache create file error: %s", error));
		os.exit(1);
	    }
	    fs.close(result[0]);

	    // Write file
	    error = fs.write(path, JSON.stringify(value), 0644);
	    if (error) {
		log(sprintf("Cache write file error: %s", error));
		os.exit(1);
	    }

	    // Write the expiration.
	    this.exMap[key] = expiry;
	    this.writeMeta();
	}
    };

    // Export
    if (typeof(exports) != 'undefined') {
	exports.Cache = Cache;
    }
})(typeof window === 'undefined' ? this : window);
