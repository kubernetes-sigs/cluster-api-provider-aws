var __webpack_modules__ = {
    761: module => {
        module.exports = r => {
            const n = process.versions.node.split(".").map((x => parseInt(x, 10)));
            r = r.split(".").map((x => parseInt(x, 10)));
            return n[0] > r[0] || n[0] === r[0] && (n[1] > r[1] || n[1] === r[1] && n[2] >= r[2]);
        };
    },
    2047: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(5747);
        const path = __webpack_require__(5622);
        const LCHOWN = fs.lchown ? "lchown" : "chown";
        const LCHOWNSYNC = fs.lchownSync ? "lchownSync" : "chownSync";
        const needEISDIRHandled = fs.lchown && !process.version.match(/v1[1-9]+\./) && !process.version.match(/v10\.[6-9]/);
        const lchownSync = (path, uid, gid) => {
            try {
                return fs[LCHOWNSYNC](path, uid, gid);
            } catch (er) {
                if (er.code !== "ENOENT") throw er;
            }
        };
        const chownSync = (path, uid, gid) => {
            try {
                return fs.chownSync(path, uid, gid);
            } catch (er) {
                if (er.code !== "ENOENT") throw er;
            }
        };
        const handleEISDIR = needEISDIRHandled ? (path, uid, gid, cb) => er => {
            if (!er || er.code !== "EISDIR") cb(er); else fs.chown(path, uid, gid, cb);
        } : (_, __, ___, cb) => cb;
        const handleEISDirSync = needEISDIRHandled ? (path, uid, gid) => {
            try {
                return lchownSync(path, uid, gid);
            } catch (er) {
                if (er.code !== "EISDIR") throw er;
                chownSync(path, uid, gid);
            }
        } : (path, uid, gid) => lchownSync(path, uid, gid);
        const nodeVersion = process.version;
        let readdir = (path, options, cb) => fs.readdir(path, options, cb);
        let readdirSync = (path, options) => fs.readdirSync(path, options);
        if (/^v4\./.test(nodeVersion)) readdir = (path, options, cb) => fs.readdir(path, cb);
        const chown = (cpath, uid, gid, cb) => {
            fs[LCHOWN](cpath, uid, gid, handleEISDIR(cpath, uid, gid, (er => {
                cb(er && er.code !== "ENOENT" ? er : null);
            })));
        };
        const chownrKid = (p, child, uid, gid, cb) => {
            if (typeof child === "string") return fs.lstat(path.resolve(p, child), ((er, stats) => {
                if (er) return cb(er.code !== "ENOENT" ? er : null);
                stats.name = child;
                chownrKid(p, stats, uid, gid, cb);
            }));
            if (child.isDirectory()) {
                chownr(path.resolve(p, child.name), uid, gid, (er => {
                    if (er) return cb(er);
                    const cpath = path.resolve(p, child.name);
                    chown(cpath, uid, gid, cb);
                }));
            } else {
                const cpath = path.resolve(p, child.name);
                chown(cpath, uid, gid, cb);
            }
        };
        const chownr = (p, uid, gid, cb) => {
            readdir(p, {
                withFileTypes: true
            }, ((er, children) => {
                if (er) {
                    if (er.code === "ENOENT") return cb(); else if (er.code !== "ENOTDIR" && er.code !== "ENOTSUP") return cb(er);
                }
                if (er || !children.length) return chown(p, uid, gid, cb);
                let len = children.length;
                let errState = null;
                const then = er => {
                    if (errState) return;
                    if (er) return cb(errState = er);
                    if (--len === 0) return chown(p, uid, gid, cb);
                };
                children.forEach((child => chownrKid(p, child, uid, gid, then)));
            }));
        };
        const chownrKidSync = (p, child, uid, gid) => {
            if (typeof child === "string") {
                try {
                    const stats = fs.lstatSync(path.resolve(p, child));
                    stats.name = child;
                    child = stats;
                } catch (er) {
                    if (er.code === "ENOENT") return; else throw er;
                }
            }
            if (child.isDirectory()) chownrSync(path.resolve(p, child.name), uid, gid);
            handleEISDirSync(path.resolve(p, child.name), uid, gid);
        };
        const chownrSync = (p, uid, gid) => {
            let children;
            try {
                children = readdirSync(p, {
                    withFileTypes: true
                });
            } catch (er) {
                if (er.code === "ENOENT") return; else if (er.code === "ENOTDIR" || er.code === "ENOTSUP") return handleEISDirSync(p, uid, gid); else throw er;
            }
            if (children && children.length) children.forEach((child => chownrKidSync(p, child, uid, gid)));
            return handleEISDirSync(p, uid, gid);
        };
        module.exports = chownr;
        chownr.sync = chownrSync;
    },
    2197: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(6851);
        const path = __webpack_require__(5622);
        const mkdirsSync = __webpack_require__(7311).mkdirsSync;
        const utimesMillisSync = __webpack_require__(5302).utimesMillisSync;
        const stat = __webpack_require__(6637);
        function copySync(src, dest, opts) {
            if (typeof opts === "function") {
                opts = {
                    filter: opts
                };
            }
            opts = opts || {};
            opts.clobber = "clobber" in opts ? !!opts.clobber : true;
            opts.overwrite = "overwrite" in opts ? !!opts.overwrite : opts.clobber;
            if (opts.preserveTimestamps && process.arch === "ia32") {
                console.warn(`fs-extra: Using the preserveTimestamps option in 32-bit node is not recommended;\n\n    see https://github.com/jprichardson/node-fs-extra/issues/269`);
            }
            const {srcStat, destStat} = stat.checkPathsSync(src, dest, "copy");
            stat.checkParentPathsSync(src, srcStat, dest, "copy");
            return handleFilterAndCopy(destStat, src, dest, opts);
        }
        function handleFilterAndCopy(destStat, src, dest, opts) {
            if (opts.filter && !opts.filter(src, dest)) return;
            const destParent = path.dirname(dest);
            if (!fs.existsSync(destParent)) mkdirsSync(destParent);
            return startCopy(destStat, src, dest, opts);
        }
        function startCopy(destStat, src, dest, opts) {
            if (opts.filter && !opts.filter(src, dest)) return;
            return getStats(destStat, src, dest, opts);
        }
        function getStats(destStat, src, dest, opts) {
            const statSync = opts.dereference ? fs.statSync : fs.lstatSync;
            const srcStat = statSync(src);
            if (srcStat.isDirectory()) return onDir(srcStat, destStat, src, dest, opts); else if (srcStat.isFile() || srcStat.isCharacterDevice() || srcStat.isBlockDevice()) return onFile(srcStat, destStat, src, dest, opts); else if (srcStat.isSymbolicLink()) return onLink(destStat, src, dest, opts);
        }
        function onFile(srcStat, destStat, src, dest, opts) {
            if (!destStat) return copyFile(srcStat, src, dest, opts);
            return mayCopyFile(srcStat, src, dest, opts);
        }
        function mayCopyFile(srcStat, src, dest, opts) {
            if (opts.overwrite) {
                fs.unlinkSync(dest);
                return copyFile(srcStat, src, dest, opts);
            } else if (opts.errorOnExist) {
                throw new Error(`'${dest}' already exists`);
            }
        }
        function copyFile(srcStat, src, dest, opts) {
            fs.copyFileSync(src, dest);
            if (opts.preserveTimestamps) handleTimestamps(srcStat.mode, src, dest);
            return setDestMode(dest, srcStat.mode);
        }
        function handleTimestamps(srcMode, src, dest) {
            if (fileIsNotWritable(srcMode)) makeFileWritable(dest, srcMode);
            return setDestTimestamps(src, dest);
        }
        function fileIsNotWritable(srcMode) {
            return (srcMode & 128) === 0;
        }
        function makeFileWritable(dest, srcMode) {
            return setDestMode(dest, srcMode | 128);
        }
        function setDestMode(dest, srcMode) {
            return fs.chmodSync(dest, srcMode);
        }
        function setDestTimestamps(src, dest) {
            const updatedSrcStat = fs.statSync(src);
            return utimesMillisSync(dest, updatedSrcStat.atime, updatedSrcStat.mtime);
        }
        function onDir(srcStat, destStat, src, dest, opts) {
            if (!destStat) return mkDirAndCopy(srcStat.mode, src, dest, opts);
            if (destStat && !destStat.isDirectory()) {
                throw new Error(`Cannot overwrite non-directory '${dest}' with directory '${src}'.`);
            }
            return copyDir(src, dest, opts);
        }
        function mkDirAndCopy(srcMode, src, dest, opts) {
            fs.mkdirSync(dest);
            copyDir(src, dest, opts);
            return setDestMode(dest, srcMode);
        }
        function copyDir(src, dest, opts) {
            fs.readdirSync(src).forEach((item => copyDirItem(item, src, dest, opts)));
        }
        function copyDirItem(item, src, dest, opts) {
            const srcItem = path.join(src, item);
            const destItem = path.join(dest, item);
            const {destStat} = stat.checkPathsSync(srcItem, destItem, "copy");
            return startCopy(destStat, srcItem, destItem, opts);
        }
        function onLink(destStat, src, dest, opts) {
            let resolvedSrc = fs.readlinkSync(src);
            if (opts.dereference) {
                resolvedSrc = path.resolve(process.cwd(), resolvedSrc);
            }
            if (!destStat) {
                return fs.symlinkSync(resolvedSrc, dest);
            } else {
                let resolvedDest;
                try {
                    resolvedDest = fs.readlinkSync(dest);
                } catch (err) {
                    if (err.code === "EINVAL" || err.code === "UNKNOWN") return fs.symlinkSync(resolvedSrc, dest);
                    throw err;
                }
                if (opts.dereference) {
                    resolvedDest = path.resolve(process.cwd(), resolvedDest);
                }
                if (stat.isSrcSubdir(resolvedSrc, resolvedDest)) {
                    throw new Error(`Cannot copy '${resolvedSrc}' to a subdirectory of itself, '${resolvedDest}'.`);
                }
                if (fs.statSync(dest).isDirectory() && stat.isSrcSubdir(resolvedDest, resolvedSrc)) {
                    throw new Error(`Cannot overwrite '${resolvedDest}' with '${resolvedSrc}'.`);
                }
                return copyLink(resolvedSrc, dest);
            }
        }
        function copyLink(resolvedSrc, dest) {
            fs.unlinkSync(dest);
            return fs.symlinkSync(resolvedSrc, dest);
        }
        module.exports = copySync;
    },
    8992: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        module.exports = {
            copySync: __webpack_require__(2197)
        };
    },
    465: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(6851);
        const path = __webpack_require__(5622);
        const mkdirs = __webpack_require__(7311).mkdirs;
        const pathExists = __webpack_require__(2569).pathExists;
        const utimesMillis = __webpack_require__(5302).utimesMillis;
        const stat = __webpack_require__(6637);
        function copy(src, dest, opts, cb) {
            if (typeof opts === "function" && !cb) {
                cb = opts;
                opts = {};
            } else if (typeof opts === "function") {
                opts = {
                    filter: opts
                };
            }
            cb = cb || function() {};
            opts = opts || {};
            opts.clobber = "clobber" in opts ? !!opts.clobber : true;
            opts.overwrite = "overwrite" in opts ? !!opts.overwrite : opts.clobber;
            if (opts.preserveTimestamps && process.arch === "ia32") {
                console.warn(`fs-extra: Using the preserveTimestamps option in 32-bit node is not recommended;\n\n    see https://github.com/jprichardson/node-fs-extra/issues/269`);
            }
            stat.checkPaths(src, dest, "copy", ((err, stats) => {
                if (err) return cb(err);
                const {srcStat, destStat} = stats;
                stat.checkParentPaths(src, srcStat, dest, "copy", (err => {
                    if (err) return cb(err);
                    if (opts.filter) return handleFilter(checkParentDir, destStat, src, dest, opts, cb);
                    return checkParentDir(destStat, src, dest, opts, cb);
                }));
            }));
        }
        function checkParentDir(destStat, src, dest, opts, cb) {
            const destParent = path.dirname(dest);
            pathExists(destParent, ((err, dirExists) => {
                if (err) return cb(err);
                if (dirExists) return startCopy(destStat, src, dest, opts, cb);
                mkdirs(destParent, (err => {
                    if (err) return cb(err);
                    return startCopy(destStat, src, dest, opts, cb);
                }));
            }));
        }
        function handleFilter(onInclude, destStat, src, dest, opts, cb) {
            Promise.resolve(opts.filter(src, dest)).then((include => {
                if (include) return onInclude(destStat, src, dest, opts, cb);
                return cb();
            }), (error => cb(error)));
        }
        function startCopy(destStat, src, dest, opts, cb) {
            if (opts.filter) return handleFilter(getStats, destStat, src, dest, opts, cb);
            return getStats(destStat, src, dest, opts, cb);
        }
        function getStats(destStat, src, dest, opts, cb) {
            const stat = opts.dereference ? fs.stat : fs.lstat;
            stat(src, ((err, srcStat) => {
                if (err) return cb(err);
                if (srcStat.isDirectory()) return onDir(srcStat, destStat, src, dest, opts, cb); else if (srcStat.isFile() || srcStat.isCharacterDevice() || srcStat.isBlockDevice()) return onFile(srcStat, destStat, src, dest, opts, cb); else if (srcStat.isSymbolicLink()) return onLink(destStat, src, dest, opts, cb);
            }));
        }
        function onFile(srcStat, destStat, src, dest, opts, cb) {
            if (!destStat) return copyFile(srcStat, src, dest, opts, cb);
            return mayCopyFile(srcStat, src, dest, opts, cb);
        }
        function mayCopyFile(srcStat, src, dest, opts, cb) {
            if (opts.overwrite) {
                fs.unlink(dest, (err => {
                    if (err) return cb(err);
                    return copyFile(srcStat, src, dest, opts, cb);
                }));
            } else if (opts.errorOnExist) {
                return cb(new Error(`'${dest}' already exists`));
            } else return cb();
        }
        function copyFile(srcStat, src, dest, opts, cb) {
            fs.copyFile(src, dest, (err => {
                if (err) return cb(err);
                if (opts.preserveTimestamps) return handleTimestampsAndMode(srcStat.mode, src, dest, cb);
                return setDestMode(dest, srcStat.mode, cb);
            }));
        }
        function handleTimestampsAndMode(srcMode, src, dest, cb) {
            if (fileIsNotWritable(srcMode)) {
                return makeFileWritable(dest, srcMode, (err => {
                    if (err) return cb(err);
                    return setDestTimestampsAndMode(srcMode, src, dest, cb);
                }));
            }
            return setDestTimestampsAndMode(srcMode, src, dest, cb);
        }
        function fileIsNotWritable(srcMode) {
            return (srcMode & 128) === 0;
        }
        function makeFileWritable(dest, srcMode, cb) {
            return setDestMode(dest, srcMode | 128, cb);
        }
        function setDestTimestampsAndMode(srcMode, src, dest, cb) {
            setDestTimestamps(src, dest, (err => {
                if (err) return cb(err);
                return setDestMode(dest, srcMode, cb);
            }));
        }
        function setDestMode(dest, srcMode, cb) {
            return fs.chmod(dest, srcMode, cb);
        }
        function setDestTimestamps(src, dest, cb) {
            fs.stat(src, ((err, updatedSrcStat) => {
                if (err) return cb(err);
                return utimesMillis(dest, updatedSrcStat.atime, updatedSrcStat.mtime, cb);
            }));
        }
        function onDir(srcStat, destStat, src, dest, opts, cb) {
            if (!destStat) return mkDirAndCopy(srcStat.mode, src, dest, opts, cb);
            if (destStat && !destStat.isDirectory()) {
                return cb(new Error(`Cannot overwrite non-directory '${dest}' with directory '${src}'.`));
            }
            return copyDir(src, dest, opts, cb);
        }
        function mkDirAndCopy(srcMode, src, dest, opts, cb) {
            fs.mkdir(dest, (err => {
                if (err) return cb(err);
                copyDir(src, dest, opts, (err => {
                    if (err) return cb(err);
                    return setDestMode(dest, srcMode, cb);
                }));
            }));
        }
        function copyDir(src, dest, opts, cb) {
            fs.readdir(src, ((err, items) => {
                if (err) return cb(err);
                return copyDirItems(items, src, dest, opts, cb);
            }));
        }
        function copyDirItems(items, src, dest, opts, cb) {
            const item = items.pop();
            if (!item) return cb();
            return copyDirItem(items, item, src, dest, opts, cb);
        }
        function copyDirItem(items, item, src, dest, opts, cb) {
            const srcItem = path.join(src, item);
            const destItem = path.join(dest, item);
            stat.checkPaths(srcItem, destItem, "copy", ((err, stats) => {
                if (err) return cb(err);
                const {destStat} = stats;
                startCopy(destStat, srcItem, destItem, opts, (err => {
                    if (err) return cb(err);
                    return copyDirItems(items, src, dest, opts, cb);
                }));
            }));
        }
        function onLink(destStat, src, dest, opts, cb) {
            fs.readlink(src, ((err, resolvedSrc) => {
                if (err) return cb(err);
                if (opts.dereference) {
                    resolvedSrc = path.resolve(process.cwd(), resolvedSrc);
                }
                if (!destStat) {
                    return fs.symlink(resolvedSrc, dest, cb);
                } else {
                    fs.readlink(dest, ((err, resolvedDest) => {
                        if (err) {
                            if (err.code === "EINVAL" || err.code === "UNKNOWN") return fs.symlink(resolvedSrc, dest, cb);
                            return cb(err);
                        }
                        if (opts.dereference) {
                            resolvedDest = path.resolve(process.cwd(), resolvedDest);
                        }
                        if (stat.isSrcSubdir(resolvedSrc, resolvedDest)) {
                            return cb(new Error(`Cannot copy '${resolvedSrc}' to a subdirectory of itself, '${resolvedDest}'.`));
                        }
                        if (destStat.isDirectory() && stat.isSrcSubdir(resolvedDest, resolvedSrc)) {
                            return cb(new Error(`Cannot overwrite '${resolvedDest}' with '${resolvedSrc}'.`));
                        }
                        return copyLink(resolvedSrc, dest, cb);
                    }));
                }
            }));
        }
        function copyLink(resolvedSrc, dest, cb) {
            fs.unlink(dest, (err => {
                if (err) return cb(err);
                return fs.symlink(resolvedSrc, dest, cb);
            }));
        }
        module.exports = copy;
    },
    6430: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        module.exports = {
            copy: u(__webpack_require__(465))
        };
    },
    801: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        const fs = __webpack_require__(6851);
        const path = __webpack_require__(5622);
        const mkdir = __webpack_require__(7311);
        const remove = __webpack_require__(9117);
        const emptyDir = u((function emptyDir(dir, callback) {
            callback = callback || function() {};
            fs.readdir(dir, ((err, items) => {
                if (err) return mkdir.mkdirs(dir, callback);
                items = items.map((item => path.join(dir, item)));
                deleteItem();
                function deleteItem() {
                    const item = items.pop();
                    if (!item) return callback();
                    remove.remove(item, (err => {
                        if (err) return callback(err);
                        deleteItem();
                    }));
                }
            }));
        }));
        function emptyDirSync(dir) {
            let items;
            try {
                items = fs.readdirSync(dir);
            } catch {
                return mkdir.mkdirsSync(dir);
            }
            items.forEach((item => {
                item = path.join(dir, item);
                remove.removeSync(item);
            }));
        }
        module.exports = {
            emptyDirSync,
            emptydirSync: emptyDirSync,
            emptyDir,
            emptydir: emptyDir
        };
    },
    7392: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        const path = __webpack_require__(5622);
        const fs = __webpack_require__(6851);
        const mkdir = __webpack_require__(7311);
        function createFile(file, callback) {
            function makeFile() {
                fs.writeFile(file, "", (err => {
                    if (err) return callback(err);
                    callback();
                }));
            }
            fs.stat(file, ((err, stats) => {
                if (!err && stats.isFile()) return callback();
                const dir = path.dirname(file);
                fs.stat(dir, ((err, stats) => {
                    if (err) {
                        if (err.code === "ENOENT") {
                            return mkdir.mkdirs(dir, (err => {
                                if (err) return callback(err);
                                makeFile();
                            }));
                        }
                        return callback(err);
                    }
                    if (stats.isDirectory()) makeFile(); else {
                        fs.readdir(dir, (err => {
                            if (err) return callback(err);
                        }));
                    }
                }));
            }));
        }
        function createFileSync(file) {
            let stats;
            try {
                stats = fs.statSync(file);
            } catch {}
            if (stats && stats.isFile()) return;
            const dir = path.dirname(file);
            try {
                if (!fs.statSync(dir).isDirectory()) {
                    fs.readdirSync(dir);
                }
            } catch (err) {
                if (err && err.code === "ENOENT") mkdir.mkdirsSync(dir); else throw err;
            }
            fs.writeFileSync(file, "");
        }
        module.exports = {
            createFile: u(createFile),
            createFileSync
        };
    },
    8985: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const file = __webpack_require__(7392);
        const link = __webpack_require__(8261);
        const symlink = __webpack_require__(7618);
        module.exports = {
            createFile: file.createFile,
            createFileSync: file.createFileSync,
            ensureFile: file.createFile,
            ensureFileSync: file.createFileSync,
            createLink: link.createLink,
            createLinkSync: link.createLinkSync,
            ensureLink: link.createLink,
            ensureLinkSync: link.createLinkSync,
            createSymlink: symlink.createSymlink,
            createSymlinkSync: symlink.createSymlinkSync,
            ensureSymlink: symlink.createSymlink,
            ensureSymlinkSync: symlink.createSymlinkSync
        };
    },
    8261: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        const path = __webpack_require__(5622);
        const fs = __webpack_require__(6851);
        const mkdir = __webpack_require__(7311);
        const pathExists = __webpack_require__(2569).pathExists;
        function createLink(srcpath, dstpath, callback) {
            function makeLink(srcpath, dstpath) {
                fs.link(srcpath, dstpath, (err => {
                    if (err) return callback(err);
                    callback(null);
                }));
            }
            pathExists(dstpath, ((err, destinationExists) => {
                if (err) return callback(err);
                if (destinationExists) return callback(null);
                fs.lstat(srcpath, (err => {
                    if (err) {
                        err.message = err.message.replace("lstat", "ensureLink");
                        return callback(err);
                    }
                    const dir = path.dirname(dstpath);
                    pathExists(dir, ((err, dirExists) => {
                        if (err) return callback(err);
                        if (dirExists) return makeLink(srcpath, dstpath);
                        mkdir.mkdirs(dir, (err => {
                            if (err) return callback(err);
                            makeLink(srcpath, dstpath);
                        }));
                    }));
                }));
            }));
        }
        function createLinkSync(srcpath, dstpath) {
            const destinationExists = fs.existsSync(dstpath);
            if (destinationExists) return undefined;
            try {
                fs.lstatSync(srcpath);
            } catch (err) {
                err.message = err.message.replace("lstat", "ensureLink");
                throw err;
            }
            const dir = path.dirname(dstpath);
            const dirExists = fs.existsSync(dir);
            if (dirExists) return fs.linkSync(srcpath, dstpath);
            mkdir.mkdirsSync(dir);
            return fs.linkSync(srcpath, dstpath);
        }
        module.exports = {
            createLink: u(createLink),
            createLinkSync
        };
    },
    1249: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const path = __webpack_require__(5622);
        const fs = __webpack_require__(6851);
        const pathExists = __webpack_require__(2569).pathExists;
        function symlinkPaths(srcpath, dstpath, callback) {
            if (path.isAbsolute(srcpath)) {
                return fs.lstat(srcpath, (err => {
                    if (err) {
                        err.message = err.message.replace("lstat", "ensureSymlink");
                        return callback(err);
                    }
                    return callback(null, {
                        toCwd: srcpath,
                        toDst: srcpath
                    });
                }));
            } else {
                const dstdir = path.dirname(dstpath);
                const relativeToDst = path.join(dstdir, srcpath);
                return pathExists(relativeToDst, ((err, exists) => {
                    if (err) return callback(err);
                    if (exists) {
                        return callback(null, {
                            toCwd: relativeToDst,
                            toDst: srcpath
                        });
                    } else {
                        return fs.lstat(srcpath, (err => {
                            if (err) {
                                err.message = err.message.replace("lstat", "ensureSymlink");
                                return callback(err);
                            }
                            return callback(null, {
                                toCwd: srcpath,
                                toDst: path.relative(dstdir, srcpath)
                            });
                        }));
                    }
                }));
            }
        }
        function symlinkPathsSync(srcpath, dstpath) {
            let exists;
            if (path.isAbsolute(srcpath)) {
                exists = fs.existsSync(srcpath);
                if (!exists) throw new Error("absolute srcpath does not exist");
                return {
                    toCwd: srcpath,
                    toDst: srcpath
                };
            } else {
                const dstdir = path.dirname(dstpath);
                const relativeToDst = path.join(dstdir, srcpath);
                exists = fs.existsSync(relativeToDst);
                if (exists) {
                    return {
                        toCwd: relativeToDst,
                        toDst: srcpath
                    };
                } else {
                    exists = fs.existsSync(srcpath);
                    if (!exists) throw new Error("relative srcpath does not exist");
                    return {
                        toCwd: srcpath,
                        toDst: path.relative(dstdir, srcpath)
                    };
                }
            }
        }
        module.exports = {
            symlinkPaths,
            symlinkPathsSync
        };
    },
    8065: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(6851);
        function symlinkType(srcpath, type, callback) {
            callback = typeof type === "function" ? type : callback;
            type = typeof type === "function" ? false : type;
            if (type) return callback(null, type);
            fs.lstat(srcpath, ((err, stats) => {
                if (err) return callback(null, "file");
                type = stats && stats.isDirectory() ? "dir" : "file";
                callback(null, type);
            }));
        }
        function symlinkTypeSync(srcpath, type) {
            let stats;
            if (type) return type;
            try {
                stats = fs.lstatSync(srcpath);
            } catch {
                return "file";
            }
            return stats && stats.isDirectory() ? "dir" : "file";
        }
        module.exports = {
            symlinkType,
            symlinkTypeSync
        };
    },
    7618: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        const path = __webpack_require__(5622);
        const fs = __webpack_require__(6851);
        const _mkdirs = __webpack_require__(7311);
        const mkdirs = _mkdirs.mkdirs;
        const mkdirsSync = _mkdirs.mkdirsSync;
        const _symlinkPaths = __webpack_require__(1249);
        const symlinkPaths = _symlinkPaths.symlinkPaths;
        const symlinkPathsSync = _symlinkPaths.symlinkPathsSync;
        const _symlinkType = __webpack_require__(8065);
        const symlinkType = _symlinkType.symlinkType;
        const symlinkTypeSync = _symlinkType.symlinkTypeSync;
        const pathExists = __webpack_require__(2569).pathExists;
        function createSymlink(srcpath, dstpath, type, callback) {
            callback = typeof type === "function" ? type : callback;
            type = typeof type === "function" ? false : type;
            pathExists(dstpath, ((err, destinationExists) => {
                if (err) return callback(err);
                if (destinationExists) return callback(null);
                symlinkPaths(srcpath, dstpath, ((err, relative) => {
                    if (err) return callback(err);
                    srcpath = relative.toDst;
                    symlinkType(relative.toCwd, type, ((err, type) => {
                        if (err) return callback(err);
                        const dir = path.dirname(dstpath);
                        pathExists(dir, ((err, dirExists) => {
                            if (err) return callback(err);
                            if (dirExists) return fs.symlink(srcpath, dstpath, type, callback);
                            mkdirs(dir, (err => {
                                if (err) return callback(err);
                                fs.symlink(srcpath, dstpath, type, callback);
                            }));
                        }));
                    }));
                }));
            }));
        }
        function createSymlinkSync(srcpath, dstpath, type) {
            const destinationExists = fs.existsSync(dstpath);
            if (destinationExists) return undefined;
            const relative = symlinkPathsSync(srcpath, dstpath);
            srcpath = relative.toDst;
            type = symlinkTypeSync(relative.toCwd, type);
            const dir = path.dirname(dstpath);
            const exists = fs.existsSync(dir);
            if (exists) return fs.symlinkSync(srcpath, dstpath, type);
            mkdirsSync(dir);
            return fs.symlinkSync(srcpath, dstpath, type);
        }
        module.exports = {
            createSymlink: u(createSymlink),
            createSymlinkSync
        };
    },
    5093: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        const fs = __webpack_require__(6851);
        const api = [ "access", "appendFile", "chmod", "chown", "close", "copyFile", "fchmod", "fchown", "fdatasync", "fstat", "fsync", "ftruncate", "futimes", "lchmod", "lchown", "link", "lstat", "mkdir", "mkdtemp", "open", "opendir", "readdir", "readFile", "readlink", "realpath", "rename", "rm", "rmdir", "stat", "symlink", "truncate", "unlink", "utimes", "writeFile" ].filter((key => typeof fs[key] === "function"));
        Object.keys(fs).forEach((key => {
            if (key === "promises") {
                return;
            }
            exports[key] = fs[key];
        }));
        api.forEach((method => {
            exports[method] = u(fs[method]);
        }));
        exports.exists = function(filename, callback) {
            if (typeof callback === "function") {
                return fs.exists(filename, callback);
            }
            return new Promise((resolve => fs.exists(filename, resolve)));
        };
        exports.read = function(fd, buffer, offset, length, position, callback) {
            if (typeof callback === "function") {
                return fs.read(fd, buffer, offset, length, position, callback);
            }
            return new Promise(((resolve, reject) => {
                fs.read(fd, buffer, offset, length, position, ((err, bytesRead, buffer) => {
                    if (err) return reject(err);
                    resolve({
                        bytesRead,
                        buffer
                    });
                }));
            }));
        };
        exports.write = function(fd, buffer, ...args) {
            if (typeof args[args.length - 1] === "function") {
                return fs.write(fd, buffer, ...args);
            }
            return new Promise(((resolve, reject) => {
                fs.write(fd, buffer, ...args, ((err, bytesWritten, buffer) => {
                    if (err) return reject(err);
                    resolve({
                        bytesWritten,
                        buffer
                    });
                }));
            }));
        };
        if (typeof fs.writev === "function") {
            exports.writev = function(fd, buffers, ...args) {
                if (typeof args[args.length - 1] === "function") {
                    return fs.writev(fd, buffers, ...args);
                }
                return new Promise(((resolve, reject) => {
                    fs.writev(fd, buffers, ...args, ((err, bytesWritten, buffers) => {
                        if (err) return reject(err);
                        resolve({
                            bytesWritten,
                            buffers
                        });
                    }));
                }));
            };
        }
        if (typeof fs.realpath.native === "function") {
            exports.realpath.native = u(fs.realpath.native);
        }
    },
    9728: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        module.exports = {
            ...__webpack_require__(5093),
            ...__webpack_require__(8992),
            ...__webpack_require__(6430),
            ...__webpack_require__(801),
            ...__webpack_require__(8985),
            ...__webpack_require__(3779),
            ...__webpack_require__(7311),
            ...__webpack_require__(3815),
            ...__webpack_require__(1034),
            ...__webpack_require__(836),
            ...__webpack_require__(2569),
            ...__webpack_require__(9117)
        };
        const fs = __webpack_require__(5747);
        if (Object.getOwnPropertyDescriptor(fs, "promises")) {
            Object.defineProperty(module.exports, "promises", {
                get() {
                    return fs.promises;
                }
            });
        }
    },
    3779: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).p;
        const jsonFile = __webpack_require__(2002);
        jsonFile.outputJson = u(__webpack_require__(209));
        jsonFile.outputJsonSync = __webpack_require__(8757);
        jsonFile.outputJSON = jsonFile.outputJson;
        jsonFile.outputJSONSync = jsonFile.outputJsonSync;
        jsonFile.writeJSON = jsonFile.writeJson;
        jsonFile.writeJSONSync = jsonFile.writeJsonSync;
        jsonFile.readJSON = jsonFile.readJson;
        jsonFile.readJSONSync = jsonFile.readJsonSync;
        module.exports = jsonFile;
    },
    2002: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const jsonFile = __webpack_require__(3393);
        module.exports = {
            readJson: jsonFile.readFile,
            readJsonSync: jsonFile.readFileSync,
            writeJson: jsonFile.writeFile,
            writeJsonSync: jsonFile.writeFileSync
        };
    },
    8757: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const {stringify} = __webpack_require__(9293);
        const {outputFileSync} = __webpack_require__(836);
        function outputJsonSync(file, data, options) {
            const str = stringify(data, options);
            outputFileSync(file, str, options);
        }
        module.exports = outputJsonSync;
    },
    209: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const {stringify} = __webpack_require__(9293);
        const {outputFile} = __webpack_require__(836);
        async function outputJson(file, data, options = {}) {
            const str = stringify(data, options);
            await outputFile(file, str, options);
        }
        module.exports = outputJson;
    },
    7311: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).p;
        const {makeDir: _makeDir, makeDirSync} = __webpack_require__(3057);
        const makeDir = u(_makeDir);
        module.exports = {
            mkdirs: makeDir,
            mkdirsSync: makeDirSync,
            mkdirp: makeDir,
            mkdirpSync: makeDirSync,
            ensureDir: makeDir,
            ensureDirSync: makeDirSync
        };
    },
    3057: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(5093);
        const path = __webpack_require__(5622);
        const atLeastNode = __webpack_require__(761);
        const useNativeRecursiveOption = atLeastNode("10.12.0");
        const checkPath = pth => {
            if (process.platform === "win32") {
                const pathHasInvalidWinCharacters = /[<>:"|?*]/.test(pth.replace(path.parse(pth).root, ""));
                if (pathHasInvalidWinCharacters) {
                    const error = new Error(`Path contains invalid characters: ${pth}`);
                    error.code = "EINVAL";
                    throw error;
                }
            }
        };
        const processOptions = options => {
            const defaults = {
                mode: 511
            };
            if (typeof options === "number") options = {
                mode: options
            };
            return {
                ...defaults,
                ...options
            };
        };
        const permissionError = pth => {
            const error = new Error(`operation not permitted, mkdir '${pth}'`);
            error.code = "EPERM";
            error.errno = -4048;
            error.path = pth;
            error.syscall = "mkdir";
            return error;
        };
        module.exports.makeDir = async (input, options) => {
            checkPath(input);
            options = processOptions(options);
            if (useNativeRecursiveOption) {
                const pth = path.resolve(input);
                return fs.mkdir(pth, {
                    mode: options.mode,
                    recursive: true
                });
            }
            const make = async pth => {
                try {
                    await fs.mkdir(pth, options.mode);
                } catch (error) {
                    if (error.code === "EPERM") {
                        throw error;
                    }
                    if (error.code === "ENOENT") {
                        if (path.dirname(pth) === pth) {
                            throw permissionError(pth);
                        }
                        if (error.message.includes("null bytes")) {
                            throw error;
                        }
                        await make(path.dirname(pth));
                        return make(pth);
                    }
                    try {
                        const stats = await fs.stat(pth);
                        if (!stats.isDirectory()) {
                            throw new Error("The path is not a directory");
                        }
                    } catch {
                        throw error;
                    }
                }
            };
            return make(path.resolve(input));
        };
        module.exports.makeDirSync = (input, options) => {
            checkPath(input);
            options = processOptions(options);
            if (useNativeRecursiveOption) {
                const pth = path.resolve(input);
                return fs.mkdirSync(pth, {
                    mode: options.mode,
                    recursive: true
                });
            }
            const make = pth => {
                try {
                    fs.mkdirSync(pth, options.mode);
                } catch (error) {
                    if (error.code === "EPERM") {
                        throw error;
                    }
                    if (error.code === "ENOENT") {
                        if (path.dirname(pth) === pth) {
                            throw permissionError(pth);
                        }
                        if (error.message.includes("null bytes")) {
                            throw error;
                        }
                        make(path.dirname(pth));
                        return make(pth);
                    }
                    try {
                        if (!fs.statSync(pth).isDirectory()) {
                            throw new Error("The path is not a directory");
                        }
                    } catch {
                        throw error;
                    }
                }
            };
            return make(path.resolve(input));
        };
    },
    3815: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        module.exports = {
            moveSync: __webpack_require__(2727)
        };
    },
    2727: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(6851);
        const path = __webpack_require__(5622);
        const copySync = __webpack_require__(8992).copySync;
        const removeSync = __webpack_require__(9117).removeSync;
        const mkdirpSync = __webpack_require__(7311).mkdirpSync;
        const stat = __webpack_require__(6637);
        function moveSync(src, dest, opts) {
            opts = opts || {};
            const overwrite = opts.overwrite || opts.clobber || false;
            const {srcStat} = stat.checkPathsSync(src, dest, "move");
            stat.checkParentPathsSync(src, srcStat, dest, "move");
            mkdirpSync(path.dirname(dest));
            return doRename(src, dest, overwrite);
        }
        function doRename(src, dest, overwrite) {
            if (overwrite) {
                removeSync(dest);
                return rename(src, dest, overwrite);
            }
            if (fs.existsSync(dest)) throw new Error("dest already exists.");
            return rename(src, dest, overwrite);
        }
        function rename(src, dest, overwrite) {
            try {
                fs.renameSync(src, dest);
            } catch (err) {
                if (err.code !== "EXDEV") throw err;
                return moveAcrossDevice(src, dest, overwrite);
            }
        }
        function moveAcrossDevice(src, dest, overwrite) {
            const opts = {
                overwrite,
                errorOnExist: true
            };
            copySync(src, dest, opts);
            return removeSync(src);
        }
        module.exports = moveSync;
    },
    1034: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        module.exports = {
            move: u(__webpack_require__(2521))
        };
    },
    2521: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(6851);
        const path = __webpack_require__(5622);
        const copy = __webpack_require__(6430).copy;
        const remove = __webpack_require__(9117).remove;
        const mkdirp = __webpack_require__(7311).mkdirp;
        const pathExists = __webpack_require__(2569).pathExists;
        const stat = __webpack_require__(6637);
        function move(src, dest, opts, cb) {
            if (typeof opts === "function") {
                cb = opts;
                opts = {};
            }
            const overwrite = opts.overwrite || opts.clobber || false;
            stat.checkPaths(src, dest, "move", ((err, stats) => {
                if (err) return cb(err);
                const {srcStat} = stats;
                stat.checkParentPaths(src, srcStat, dest, "move", (err => {
                    if (err) return cb(err);
                    mkdirp(path.dirname(dest), (err => {
                        if (err) return cb(err);
                        return doRename(src, dest, overwrite, cb);
                    }));
                }));
            }));
        }
        function doRename(src, dest, overwrite, cb) {
            if (overwrite) {
                return remove(dest, (err => {
                    if (err) return cb(err);
                    return rename(src, dest, overwrite, cb);
                }));
            }
            pathExists(dest, ((err, destExists) => {
                if (err) return cb(err);
                if (destExists) return cb(new Error("dest already exists."));
                return rename(src, dest, overwrite, cb);
            }));
        }
        function rename(src, dest, overwrite, cb) {
            fs.rename(src, dest, (err => {
                if (!err) return cb();
                if (err.code !== "EXDEV") return cb(err);
                return moveAcrossDevice(src, dest, overwrite, cb);
            }));
        }
        function moveAcrossDevice(src, dest, overwrite, cb) {
            const opts = {
                overwrite,
                errorOnExist: true
            };
            copy(src, dest, opts, (err => {
                if (err) return cb(err);
                return remove(src, cb);
            }));
        }
        module.exports = move;
    },
    836: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        const fs = __webpack_require__(6851);
        const path = __webpack_require__(5622);
        const mkdir = __webpack_require__(7311);
        const pathExists = __webpack_require__(2569).pathExists;
        function outputFile(file, data, encoding, callback) {
            if (typeof encoding === "function") {
                callback = encoding;
                encoding = "utf8";
            }
            const dir = path.dirname(file);
            pathExists(dir, ((err, itDoes) => {
                if (err) return callback(err);
                if (itDoes) return fs.writeFile(file, data, encoding, callback);
                mkdir.mkdirs(dir, (err => {
                    if (err) return callback(err);
                    fs.writeFile(file, data, encoding, callback);
                }));
            }));
        }
        function outputFileSync(file, ...args) {
            const dir = path.dirname(file);
            if (fs.existsSync(dir)) {
                return fs.writeFileSync(file, ...args);
            }
            mkdir.mkdirsSync(dir);
            fs.writeFileSync(file, ...args);
        }
        module.exports = {
            outputFile: u(outputFile),
            outputFileSync
        };
    },
    2569: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).p;
        const fs = __webpack_require__(5093);
        function pathExists(path) {
            return fs.access(path).then((() => true)).catch((() => false));
        }
        module.exports = {
            pathExists: u(pathExists),
            pathExistsSync: fs.existsSync
        };
    },
    9117: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const u = __webpack_require__(2141).E;
        const rimraf = __webpack_require__(1683);
        module.exports = {
            remove: u(rimraf),
            removeSync: rimraf.sync
        };
    },
    1683: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(6851);
        const path = __webpack_require__(5622);
        const assert = __webpack_require__(2357);
        const isWindows = process.platform === "win32";
        function defaults(options) {
            const methods = [ "unlink", "chmod", "stat", "lstat", "rmdir", "readdir" ];
            methods.forEach((m => {
                options[m] = options[m] || fs[m];
                m = m + "Sync";
                options[m] = options[m] || fs[m];
            }));
            options.maxBusyTries = options.maxBusyTries || 3;
        }
        function rimraf(p, options, cb) {
            let busyTries = 0;
            if (typeof options === "function") {
                cb = options;
                options = {};
            }
            assert(p, "rimraf: missing path");
            assert.strictEqual(typeof p, "string", "rimraf: path should be a string");
            assert.strictEqual(typeof cb, "function", "rimraf: callback function required");
            assert(options, "rimraf: invalid options argument provided");
            assert.strictEqual(typeof options, "object", "rimraf: options should be object");
            defaults(options);
            rimraf_(p, options, (function CB(er) {
                if (er) {
                    if ((er.code === "EBUSY" || er.code === "ENOTEMPTY" || er.code === "EPERM") && busyTries < options.maxBusyTries) {
                        busyTries++;
                        const time = busyTries * 100;
                        return setTimeout((() => rimraf_(p, options, CB)), time);
                    }
                    if (er.code === "ENOENT") er = null;
                }
                cb(er);
            }));
        }
        function rimraf_(p, options, cb) {
            assert(p);
            assert(options);
            assert(typeof cb === "function");
            options.lstat(p, ((er, st) => {
                if (er && er.code === "ENOENT") {
                    return cb(null);
                }
                if (er && er.code === "EPERM" && isWindows) {
                    return fixWinEPERM(p, options, er, cb);
                }
                if (st && st.isDirectory()) {
                    return rmdir(p, options, er, cb);
                }
                options.unlink(p, (er => {
                    if (er) {
                        if (er.code === "ENOENT") {
                            return cb(null);
                        }
                        if (er.code === "EPERM") {
                            return isWindows ? fixWinEPERM(p, options, er, cb) : rmdir(p, options, er, cb);
                        }
                        if (er.code === "EISDIR") {
                            return rmdir(p, options, er, cb);
                        }
                    }
                    return cb(er);
                }));
            }));
        }
        function fixWinEPERM(p, options, er, cb) {
            assert(p);
            assert(options);
            assert(typeof cb === "function");
            options.chmod(p, 438, (er2 => {
                if (er2) {
                    cb(er2.code === "ENOENT" ? null : er);
                } else {
                    options.stat(p, ((er3, stats) => {
                        if (er3) {
                            cb(er3.code === "ENOENT" ? null : er);
                        } else if (stats.isDirectory()) {
                            rmdir(p, options, er, cb);
                        } else {
                            options.unlink(p, cb);
                        }
                    }));
                }
            }));
        }
        function fixWinEPERMSync(p, options, er) {
            let stats;
            assert(p);
            assert(options);
            try {
                options.chmodSync(p, 438);
            } catch (er2) {
                if (er2.code === "ENOENT") {
                    return;
                } else {
                    throw er;
                }
            }
            try {
                stats = options.statSync(p);
            } catch (er3) {
                if (er3.code === "ENOENT") {
                    return;
                } else {
                    throw er;
                }
            }
            if (stats.isDirectory()) {
                rmdirSync(p, options, er);
            } else {
                options.unlinkSync(p);
            }
        }
        function rmdir(p, options, originalEr, cb) {
            assert(p);
            assert(options);
            assert(typeof cb === "function");
            options.rmdir(p, (er => {
                if (er && (er.code === "ENOTEMPTY" || er.code === "EEXIST" || er.code === "EPERM")) {
                    rmkids(p, options, cb);
                } else if (er && er.code === "ENOTDIR") {
                    cb(originalEr);
                } else {
                    cb(er);
                }
            }));
        }
        function rmkids(p, options, cb) {
            assert(p);
            assert(options);
            assert(typeof cb === "function");
            options.readdir(p, ((er, files) => {
                if (er) return cb(er);
                let n = files.length;
                let errState;
                if (n === 0) return options.rmdir(p, cb);
                files.forEach((f => {
                    rimraf(path.join(p, f), options, (er => {
                        if (errState) {
                            return;
                        }
                        if (er) return cb(errState = er);
                        if (--n === 0) {
                            options.rmdir(p, cb);
                        }
                    }));
                }));
            }));
        }
        function rimrafSync(p, options) {
            let st;
            options = options || {};
            defaults(options);
            assert(p, "rimraf: missing path");
            assert.strictEqual(typeof p, "string", "rimraf: path should be a string");
            assert(options, "rimraf: missing options");
            assert.strictEqual(typeof options, "object", "rimraf: options should be object");
            try {
                st = options.lstatSync(p);
            } catch (er) {
                if (er.code === "ENOENT") {
                    return;
                }
                if (er.code === "EPERM" && isWindows) {
                    fixWinEPERMSync(p, options, er);
                }
            }
            try {
                if (st && st.isDirectory()) {
                    rmdirSync(p, options, null);
                } else {
                    options.unlinkSync(p);
                }
            } catch (er) {
                if (er.code === "ENOENT") {
                    return;
                } else if (er.code === "EPERM") {
                    return isWindows ? fixWinEPERMSync(p, options, er) : rmdirSync(p, options, er);
                } else if (er.code !== "EISDIR") {
                    throw er;
                }
                rmdirSync(p, options, er);
            }
        }
        function rmdirSync(p, options, originalEr) {
            assert(p);
            assert(options);
            try {
                options.rmdirSync(p);
            } catch (er) {
                if (er.code === "ENOTDIR") {
                    throw originalEr;
                } else if (er.code === "ENOTEMPTY" || er.code === "EEXIST" || er.code === "EPERM") {
                    rmkidsSync(p, options);
                } else if (er.code !== "ENOENT") {
                    throw er;
                }
            }
        }
        function rmkidsSync(p, options) {
            assert(p);
            assert(options);
            options.readdirSync(p).forEach((f => rimrafSync(path.join(p, f), options)));
            if (isWindows) {
                const startTime = Date.now();
                do {
                    try {
                        const ret = options.rmdirSync(p, options);
                        return ret;
                    } catch {}
                } while (Date.now() - startTime < 500);
            } else {
                const ret = options.rmdirSync(p, options);
                return ret;
            }
        }
        module.exports = rimraf;
        rimraf.sync = rimrafSync;
    },
    6637: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(5093);
        const path = __webpack_require__(5622);
        const util = __webpack_require__(1669);
        const atLeastNode = __webpack_require__(761);
        const nodeSupportsBigInt = atLeastNode("10.5.0");
        const stat = file => nodeSupportsBigInt ? fs.stat(file, {
            bigint: true
        }) : fs.stat(file);
        const statSync = file => nodeSupportsBigInt ? fs.statSync(file, {
            bigint: true
        }) : fs.statSync(file);
        function getStats(src, dest) {
            return Promise.all([ stat(src), stat(dest).catch((err => {
                if (err.code === "ENOENT") return null;
                throw err;
            })) ]).then((([srcStat, destStat]) => ({
                srcStat,
                destStat
            })));
        }
        function getStatsSync(src, dest) {
            let destStat;
            const srcStat = statSync(src);
            try {
                destStat = statSync(dest);
            } catch (err) {
                if (err.code === "ENOENT") return {
                    srcStat,
                    destStat: null
                };
                throw err;
            }
            return {
                srcStat,
                destStat
            };
        }
        function checkPaths(src, dest, funcName, cb) {
            util.callbackify(getStats)(src, dest, ((err, stats) => {
                if (err) return cb(err);
                const {srcStat, destStat} = stats;
                if (destStat && areIdentical(srcStat, destStat)) {
                    return cb(new Error("Source and destination must not be the same."));
                }
                if (srcStat.isDirectory() && isSrcSubdir(src, dest)) {
                    return cb(new Error(errMsg(src, dest, funcName)));
                }
                return cb(null, {
                    srcStat,
                    destStat
                });
            }));
        }
        function checkPathsSync(src, dest, funcName) {
            const {srcStat, destStat} = getStatsSync(src, dest);
            if (destStat && areIdentical(srcStat, destStat)) {
                throw new Error("Source and destination must not be the same.");
            }
            if (srcStat.isDirectory() && isSrcSubdir(src, dest)) {
                throw new Error(errMsg(src, dest, funcName));
            }
            return {
                srcStat,
                destStat
            };
        }
        function checkParentPaths(src, srcStat, dest, funcName, cb) {
            const srcParent = path.resolve(path.dirname(src));
            const destParent = path.resolve(path.dirname(dest));
            if (destParent === srcParent || destParent === path.parse(destParent).root) return cb();
            const callback = (err, destStat) => {
                if (err) {
                    if (err.code === "ENOENT") return cb();
                    return cb(err);
                }
                if (areIdentical(srcStat, destStat)) {
                    return cb(new Error(errMsg(src, dest, funcName)));
                }
                return checkParentPaths(src, srcStat, destParent, funcName, cb);
            };
            if (nodeSupportsBigInt) fs.stat(destParent, {
                bigint: true
            }, callback); else fs.stat(destParent, callback);
        }
        function checkParentPathsSync(src, srcStat, dest, funcName) {
            const srcParent = path.resolve(path.dirname(src));
            const destParent = path.resolve(path.dirname(dest));
            if (destParent === srcParent || destParent === path.parse(destParent).root) return;
            let destStat;
            try {
                destStat = statSync(destParent);
            } catch (err) {
                if (err.code === "ENOENT") return;
                throw err;
            }
            if (areIdentical(srcStat, destStat)) {
                throw new Error(errMsg(src, dest, funcName));
            }
            return checkParentPathsSync(src, srcStat, destParent, funcName);
        }
        function areIdentical(srcStat, destStat) {
            if (destStat.ino && destStat.dev && destStat.ino === srcStat.ino && destStat.dev === srcStat.dev) {
                if (nodeSupportsBigInt || destStat.ino < Number.MAX_SAFE_INTEGER) {
                    return true;
                }
                if (destStat.size === srcStat.size && destStat.mode === srcStat.mode && destStat.nlink === srcStat.nlink && destStat.atimeMs === srcStat.atimeMs && destStat.mtimeMs === srcStat.mtimeMs && destStat.ctimeMs === srcStat.ctimeMs && destStat.birthtimeMs === srcStat.birthtimeMs) {
                    return true;
                }
            }
            return false;
        }
        function isSrcSubdir(src, dest) {
            const srcArr = path.resolve(src).split(path.sep).filter((i => i));
            const destArr = path.resolve(dest).split(path.sep).filter((i => i));
            return srcArr.reduce(((acc, cur, i) => acc && destArr[i] === cur), true);
        }
        function errMsg(src, dest, funcName) {
            return `Cannot ${funcName} '${src}' to a subdirectory of itself, '${dest}'.`;
        }
        module.exports = {
            checkPaths,
            checkPathsSync,
            checkParentPaths,
            checkParentPathsSync,
            isSrcSubdir
        };
    },
    5302: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const fs = __webpack_require__(6851);
        function utimesMillis(path, atime, mtime, callback) {
            fs.open(path, "r+", ((err, fd) => {
                if (err) return callback(err);
                fs.futimes(fd, atime, mtime, (futimesErr => {
                    fs.close(fd, (closeErr => {
                        if (callback) callback(futimesErr || closeErr);
                    }));
                }));
            }));
        }
        function utimesMillisSync(path, atime, mtime) {
            const fd = fs.openSync(path, "r+");
            fs.futimesSync(fd, atime, mtime);
            return fs.closeSync(fd);
        }
        module.exports = {
            utimesMillis,
            utimesMillisSync
        };
    },
    2141: (__unused_webpack_module, exports) => {
        "use strict";
        exports.E = function(fn) {
            return Object.defineProperty((function(...args) {
                if (typeof args[args.length - 1] === "function") fn.apply(this, args); else {
                    return new Promise(((resolve, reject) => {
                        fn.call(this, ...args, ((err, res) => err != null ? reject(err) : resolve(res)));
                    }));
                }
            }), "name", {
                value: fn.name
            });
        };
        exports.p = function(fn) {
            return Object.defineProperty((function(...args) {
                const cb = args[args.length - 1];
                if (typeof cb !== "function") return fn.apply(this, args); else fn.apply(this, args.slice(0, -1)).then((r => cb(null, r)), cb);
            }), "name", {
                value: fn.name
            });
        };
    },
    8553: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        const MiniPass = __webpack_require__(2253);
        const EE = __webpack_require__(8614).EventEmitter;
        const fs = __webpack_require__(5747);
        let writev = fs.writev;
        if (!writev) {
            const binding = process.binding("fs");
            const FSReqWrap = binding.FSReqWrap || binding.FSReqCallback;
            writev = (fd, iovec, pos, cb) => {
                const done = (er, bw) => cb(er, bw, iovec);
                const req = new FSReqWrap;
                req.oncomplete = done;
                binding.writeBuffers(fd, iovec, pos, req);
            };
        }
        const _autoClose = Symbol("_autoClose");
        const _close = Symbol("_close");
        const _ended = Symbol("_ended");
        const _fd = Symbol("_fd");
        const _finished = Symbol("_finished");
        const _flags = Symbol("_flags");
        const _flush = Symbol("_flush");
        const _handleChunk = Symbol("_handleChunk");
        const _makeBuf = Symbol("_makeBuf");
        const _mode = Symbol("_mode");
        const _needDrain = Symbol("_needDrain");
        const _onerror = Symbol("_onerror");
        const _onopen = Symbol("_onopen");
        const _onread = Symbol("_onread");
        const _onwrite = Symbol("_onwrite");
        const _open = Symbol("_open");
        const _path = Symbol("_path");
        const _pos = Symbol("_pos");
        const _queue = Symbol("_queue");
        const _read = Symbol("_read");
        const _readSize = Symbol("_readSize");
        const _reading = Symbol("_reading");
        const _remain = Symbol("_remain");
        const _size = Symbol("_size");
        const _write = Symbol("_write");
        const _writing = Symbol("_writing");
        const _defaultFlag = Symbol("_defaultFlag");
        const _errored = Symbol("_errored");
        class ReadStream extends MiniPass {
            constructor(path, opt) {
                opt = opt || {};
                super(opt);
                this.readable = true;
                this.writable = false;
                if (typeof path !== "string") throw new TypeError("path must be a string");
                this[_errored] = false;
                this[_fd] = typeof opt.fd === "number" ? opt.fd : null;
                this[_path] = path;
                this[_readSize] = opt.readSize || 16 * 1024 * 1024;
                this[_reading] = false;
                this[_size] = typeof opt.size === "number" ? opt.size : Infinity;
                this[_remain] = this[_size];
                this[_autoClose] = typeof opt.autoClose === "boolean" ? opt.autoClose : true;
                if (typeof this[_fd] === "number") this[_read](); else this[_open]();
            }
            get fd() {
                return this[_fd];
            }
            get path() {
                return this[_path];
            }
            write() {
                throw new TypeError("this is a readable stream");
            }
            end() {
                throw new TypeError("this is a readable stream");
            }
            [_open]() {
                fs.open(this[_path], "r", ((er, fd) => this[_onopen](er, fd)));
            }
            [_onopen](er, fd) {
                if (er) this[_onerror](er); else {
                    this[_fd] = fd;
                    this.emit("open", fd);
                    this[_read]();
                }
            }
            [_makeBuf]() {
                return Buffer.allocUnsafe(Math.min(this[_readSize], this[_remain]));
            }
            [_read]() {
                if (!this[_reading]) {
                    this[_reading] = true;
                    const buf = this[_makeBuf]();
                    if (buf.length === 0) return process.nextTick((() => this[_onread](null, 0, buf)));
                    fs.read(this[_fd], buf, 0, buf.length, null, ((er, br, buf) => this[_onread](er, br, buf)));
                }
            }
            [_onread](er, br, buf) {
                this[_reading] = false;
                if (er) this[_onerror](er); else if (this[_handleChunk](br, buf)) this[_read]();
            }
            [_close]() {
                if (this[_autoClose] && typeof this[_fd] === "number") {
                    const fd = this[_fd];
                    this[_fd] = null;
                    fs.close(fd, (er => er ? this.emit("error", er) : this.emit("close")));
                }
            }
            [_onerror](er) {
                this[_reading] = true;
                this[_close]();
                this.emit("error", er);
            }
            [_handleChunk](br, buf) {
                let ret = false;
                this[_remain] -= br;
                if (br > 0) ret = super.write(br < buf.length ? buf.slice(0, br) : buf);
                if (br === 0 || this[_remain] <= 0) {
                    ret = false;
                    this[_close]();
                    super.end();
                }
                return ret;
            }
            emit(ev, data) {
                switch (ev) {
                  case "prefinish":
                  case "finish":
                    break;

                  case "drain":
                    if (typeof this[_fd] === "number") this[_read]();
                    break;

                  case "error":
                    if (this[_errored]) return;
                    this[_errored] = true;
                    return super.emit(ev, data);

                  default:
                    return super.emit(ev, data);
                }
            }
        }
        class ReadStreamSync extends ReadStream {
            [_open]() {
                let threw = true;
                try {
                    this[_onopen](null, fs.openSync(this[_path], "r"));
                    threw = false;
                } finally {
                    if (threw) this[_close]();
                }
            }
            [_read]() {
                let threw = true;
                try {
                    if (!this[_reading]) {
                        this[_reading] = true;
                        do {
                            const buf = this[_makeBuf]();
                            const br = buf.length === 0 ? 0 : fs.readSync(this[_fd], buf, 0, buf.length, null);
                            if (!this[_handleChunk](br, buf)) break;
                        } while (true);
                        this[_reading] = false;
                    }
                    threw = false;
                } finally {
                    if (threw) this[_close]();
                }
            }
            [_close]() {
                if (this[_autoClose] && typeof this[_fd] === "number") {
                    const fd = this[_fd];
                    this[_fd] = null;
                    fs.closeSync(fd);
                    this.emit("close");
                }
            }
        }
        class WriteStream extends EE {
            constructor(path, opt) {
                opt = opt || {};
                super(opt);
                this.readable = false;
                this.writable = true;
                this[_errored] = false;
                this[_writing] = false;
                this[_ended] = false;
                this[_needDrain] = false;
                this[_queue] = [];
                this[_path] = path;
                this[_fd] = typeof opt.fd === "number" ? opt.fd : null;
                this[_mode] = opt.mode === undefined ? 438 : opt.mode;
                this[_pos] = typeof opt.start === "number" ? opt.start : null;
                this[_autoClose] = typeof opt.autoClose === "boolean" ? opt.autoClose : true;
                const defaultFlag = this[_pos] !== null ? "r+" : "w";
                this[_defaultFlag] = opt.flags === undefined;
                this[_flags] = this[_defaultFlag] ? defaultFlag : opt.flags;
                if (this[_fd] === null) this[_open]();
            }
            emit(ev, data) {
                if (ev === "error") {
                    if (this[_errored]) return;
                    this[_errored] = true;
                }
                return super.emit(ev, data);
            }
            get fd() {
                return this[_fd];
            }
            get path() {
                return this[_path];
            }
            [_onerror](er) {
                this[_close]();
                this[_writing] = true;
                this.emit("error", er);
            }
            [_open]() {
                fs.open(this[_path], this[_flags], this[_mode], ((er, fd) => this[_onopen](er, fd)));
            }
            [_onopen](er, fd) {
                if (this[_defaultFlag] && this[_flags] === "r+" && er && er.code === "ENOENT") {
                    this[_flags] = "w";
                    this[_open]();
                } else if (er) this[_onerror](er); else {
                    this[_fd] = fd;
                    this.emit("open", fd);
                    this[_flush]();
                }
            }
            end(buf, enc) {
                if (buf) this.write(buf, enc);
                this[_ended] = true;
                if (!this[_writing] && !this[_queue].length && typeof this[_fd] === "number") this[_onwrite](null, 0);
                return this;
            }
            write(buf, enc) {
                if (typeof buf === "string") buf = Buffer.from(buf, enc);
                if (this[_ended]) {
                    this.emit("error", new Error("write() after end()"));
                    return false;
                }
                if (this[_fd] === null || this[_writing] || this[_queue].length) {
                    this[_queue].push(buf);
                    this[_needDrain] = true;
                    return false;
                }
                this[_writing] = true;
                this[_write](buf);
                return true;
            }
            [_write](buf) {
                fs.write(this[_fd], buf, 0, buf.length, this[_pos], ((er, bw) => this[_onwrite](er, bw)));
            }
            [_onwrite](er, bw) {
                if (er) this[_onerror](er); else {
                    if (this[_pos] !== null) this[_pos] += bw;
                    if (this[_queue].length) this[_flush](); else {
                        this[_writing] = false;
                        if (this[_ended] && !this[_finished]) {
                            this[_finished] = true;
                            this[_close]();
                            this.emit("finish");
                        } else if (this[_needDrain]) {
                            this[_needDrain] = false;
                            this.emit("drain");
                        }
                    }
                }
            }
            [_flush]() {
                if (this[_queue].length === 0) {
                    if (this[_ended]) this[_onwrite](null, 0);
                } else if (this[_queue].length === 1) this[_write](this[_queue].pop()); else {
                    const iovec = this[_queue];
                    this[_queue] = [];
                    writev(this[_fd], iovec, this[_pos], ((er, bw) => this[_onwrite](er, bw)));
                }
            }
            [_close]() {
                if (this[_autoClose] && typeof this[_fd] === "number") {
                    const fd = this[_fd];
                    this[_fd] = null;
                    fs.close(fd, (er => er ? this.emit("error", er) : this.emit("close")));
                }
            }
        }
        class WriteStreamSync extends WriteStream {
            [_open]() {
                let fd;
                if (this[_defaultFlag] && this[_flags] === "r+") {
                    try {
                        fd = fs.openSync(this[_path], this[_flags], this[_mode]);
                    } catch (er) {
                        if (er.code === "ENOENT") {
                            this[_flags] = "w";
                            return this[_open]();
                        } else throw er;
                    }
                } else fd = fs.openSync(this[_path], this[_flags], this[_mode]);
                this[_onopen](null, fd);
            }
            [_close]() {
                if (this[_autoClose] && typeof this[_fd] === "number") {
                    const fd = this[_fd];
                    this[_fd] = null;
                    fs.closeSync(fd);
                    this.emit("close");
                }
            }
            [_write](buf) {
                let threw = true;
                try {
                    this[_onwrite](null, fs.writeSync(this[_fd], buf, 0, buf.length, this[_pos]));
                    threw = false;
                } finally {
                    if (threw) try {
                        this[_close]();
                    } catch (_) {}
                }
            }
        }
        exports.ReadStream = ReadStream;
        exports.ReadStreamSync = ReadStreamSync;
        exports.WriteStream = WriteStream;
        exports.WriteStreamSync = WriteStreamSync;
    },
    9788: module => {
        "use strict";
        module.exports = clone;
        var getPrototypeOf = Object.getPrototypeOf || function(obj) {
            return obj.__proto__;
        };
        function clone(obj) {
            if (obj === null || typeof obj !== "object") return obj;
            if (obj instanceof Object) var copy = {
                __proto__: getPrototypeOf(obj)
            }; else var copy = Object.create(null);
            Object.getOwnPropertyNames(obj).forEach((function(key) {
                Object.defineProperty(copy, key, Object.getOwnPropertyDescriptor(obj, key));
            }));
            return copy;
        }
    },
    6851: (module, __unused_webpack_exports, __webpack_require__) => {
        var fs = __webpack_require__(5747);
        var polyfills = __webpack_require__(7994);
        var legacy = __webpack_require__(7885);
        var clone = __webpack_require__(9788);
        var util = __webpack_require__(1669);
        var gracefulQueue;
        var previousSymbol;
        if (typeof Symbol === "function" && typeof Symbol.for === "function") {
            gracefulQueue = Symbol.for("graceful-fs.queue");
            previousSymbol = Symbol.for("graceful-fs.previous");
        } else {
            gracefulQueue = "___graceful-fs.queue";
            previousSymbol = "___graceful-fs.previous";
        }
        function noop() {}
        function publishQueue(context, queue) {
            Object.defineProperty(context, gracefulQueue, {
                get: function() {
                    return queue;
                }
            });
        }
        var debug = noop;
        if (util.debuglog) debug = util.debuglog("gfs4"); else if (/\bgfs4\b/i.test(process.env.NODE_DEBUG || "")) debug = function() {
            var m = util.format.apply(util, arguments);
            m = "GFS4: " + m.split(/\n/).join("\nGFS4: ");
            console.error(m);
        };
        if (!fs[gracefulQueue]) {
            var queue = global[gracefulQueue] || [];
            publishQueue(fs, queue);
            fs.close = function(fs$close) {
                function close(fd, cb) {
                    return fs$close.call(fs, fd, (function(err) {
                        if (!err) {
                            retry();
                        }
                        if (typeof cb === "function") cb.apply(this, arguments);
                    }));
                }
                Object.defineProperty(close, previousSymbol, {
                    value: fs$close
                });
                return close;
            }(fs.close);
            fs.closeSync = function(fs$closeSync) {
                function closeSync(fd) {
                    fs$closeSync.apply(fs, arguments);
                    retry();
                }
                Object.defineProperty(closeSync, previousSymbol, {
                    value: fs$closeSync
                });
                return closeSync;
            }(fs.closeSync);
            if (/\bgfs4\b/i.test(process.env.NODE_DEBUG || "")) {
                process.on("exit", (function() {
                    debug(fs[gracefulQueue]);
                    __webpack_require__(2357).equal(fs[gracefulQueue].length, 0);
                }));
            }
        }
        if (!global[gracefulQueue]) {
            publishQueue(global, fs[gracefulQueue]);
        }
        module.exports = patch(clone(fs));
        if (process.env.TEST_GRACEFUL_FS_GLOBAL_PATCH && !fs.__patched) {
            module.exports = patch(fs);
            fs.__patched = true;
        }
        function patch(fs) {
            polyfills(fs);
            fs.gracefulify = patch;
            fs.createReadStream = createReadStream;
            fs.createWriteStream = createWriteStream;
            var fs$readFile = fs.readFile;
            fs.readFile = readFile;
            function readFile(path, options, cb) {
                if (typeof options === "function") cb = options, options = null;
                return go$readFile(path, options, cb);
                function go$readFile(path, options, cb) {
                    return fs$readFile(path, options, (function(err) {
                        if (err && (err.code === "EMFILE" || err.code === "ENFILE")) enqueue([ go$readFile, [ path, options, cb ] ]); else {
                            if (typeof cb === "function") cb.apply(this, arguments);
                            retry();
                        }
                    }));
                }
            }
            var fs$writeFile = fs.writeFile;
            fs.writeFile = writeFile;
            function writeFile(path, data, options, cb) {
                if (typeof options === "function") cb = options, options = null;
                return go$writeFile(path, data, options, cb);
                function go$writeFile(path, data, options, cb) {
                    return fs$writeFile(path, data, options, (function(err) {
                        if (err && (err.code === "EMFILE" || err.code === "ENFILE")) enqueue([ go$writeFile, [ path, data, options, cb ] ]); else {
                            if (typeof cb === "function") cb.apply(this, arguments);
                            retry();
                        }
                    }));
                }
            }
            var fs$appendFile = fs.appendFile;
            if (fs$appendFile) fs.appendFile = appendFile;
            function appendFile(path, data, options, cb) {
                if (typeof options === "function") cb = options, options = null;
                return go$appendFile(path, data, options, cb);
                function go$appendFile(path, data, options, cb) {
                    return fs$appendFile(path, data, options, (function(err) {
                        if (err && (err.code === "EMFILE" || err.code === "ENFILE")) enqueue([ go$appendFile, [ path, data, options, cb ] ]); else {
                            if (typeof cb === "function") cb.apply(this, arguments);
                            retry();
                        }
                    }));
                }
            }
            var fs$copyFile = fs.copyFile;
            if (fs$copyFile) fs.copyFile = copyFile;
            function copyFile(src, dest, flags, cb) {
                if (typeof flags === "function") {
                    cb = flags;
                    flags = 0;
                }
                return fs$copyFile(src, dest, flags, (function(err) {
                    if (err && (err.code === "EMFILE" || err.code === "ENFILE")) enqueue([ fs$copyFile, [ src, dest, flags, cb ] ]); else {
                        if (typeof cb === "function") cb.apply(this, arguments);
                        retry();
                    }
                }));
            }
            var fs$readdir = fs.readdir;
            fs.readdir = readdir;
            function readdir(path, options, cb) {
                var args = [ path ];
                if (typeof options !== "function") {
                    args.push(options);
                } else {
                    cb = options;
                }
                args.push(go$readdir$cb);
                return go$readdir(args);
                function go$readdir$cb(err, files) {
                    if (files && files.sort) files.sort();
                    if (err && (err.code === "EMFILE" || err.code === "ENFILE")) enqueue([ go$readdir, [ args ] ]); else {
                        if (typeof cb === "function") cb.apply(this, arguments);
                        retry();
                    }
                }
            }
            function go$readdir(args) {
                return fs$readdir.apply(fs, args);
            }
            if (process.version.substr(0, 4) === "v0.8") {
                var legStreams = legacy(fs);
                ReadStream = legStreams.ReadStream;
                WriteStream = legStreams.WriteStream;
            }
            var fs$ReadStream = fs.ReadStream;
            if (fs$ReadStream) {
                ReadStream.prototype = Object.create(fs$ReadStream.prototype);
                ReadStream.prototype.open = ReadStream$open;
            }
            var fs$WriteStream = fs.WriteStream;
            if (fs$WriteStream) {
                WriteStream.prototype = Object.create(fs$WriteStream.prototype);
                WriteStream.prototype.open = WriteStream$open;
            }
            Object.defineProperty(fs, "ReadStream", {
                get: function() {
                    return ReadStream;
                },
                set: function(val) {
                    ReadStream = val;
                },
                enumerable: true,
                configurable: true
            });
            Object.defineProperty(fs, "WriteStream", {
                get: function() {
                    return WriteStream;
                },
                set: function(val) {
                    WriteStream = val;
                },
                enumerable: true,
                configurable: true
            });
            var FileReadStream = ReadStream;
            Object.defineProperty(fs, "FileReadStream", {
                get: function() {
                    return FileReadStream;
                },
                set: function(val) {
                    FileReadStream = val;
                },
                enumerable: true,
                configurable: true
            });
            var FileWriteStream = WriteStream;
            Object.defineProperty(fs, "FileWriteStream", {
                get: function() {
                    return FileWriteStream;
                },
                set: function(val) {
                    FileWriteStream = val;
                },
                enumerable: true,
                configurable: true
            });
            function ReadStream(path, options) {
                if (this instanceof ReadStream) return fs$ReadStream.apply(this, arguments), this; else return ReadStream.apply(Object.create(ReadStream.prototype), arguments);
            }
            function ReadStream$open() {
                var that = this;
                open(that.path, that.flags, that.mode, (function(err, fd) {
                    if (err) {
                        if (that.autoClose) that.destroy();
                        that.emit("error", err);
                    } else {
                        that.fd = fd;
                        that.emit("open", fd);
                        that.read();
                    }
                }));
            }
            function WriteStream(path, options) {
                if (this instanceof WriteStream) return fs$WriteStream.apply(this, arguments), this; else return WriteStream.apply(Object.create(WriteStream.prototype), arguments);
            }
            function WriteStream$open() {
                var that = this;
                open(that.path, that.flags, that.mode, (function(err, fd) {
                    if (err) {
                        that.destroy();
                        that.emit("error", err);
                    } else {
                        that.fd = fd;
                        that.emit("open", fd);
                    }
                }));
            }
            function createReadStream(path, options) {
                return new fs.ReadStream(path, options);
            }
            function createWriteStream(path, options) {
                return new fs.WriteStream(path, options);
            }
            var fs$open = fs.open;
            fs.open = open;
            function open(path, flags, mode, cb) {
                if (typeof mode === "function") cb = mode, mode = null;
                return go$open(path, flags, mode, cb);
                function go$open(path, flags, mode, cb) {
                    return fs$open(path, flags, mode, (function(err, fd) {
                        if (err && (err.code === "EMFILE" || err.code === "ENFILE")) enqueue([ go$open, [ path, flags, mode, cb ] ]); else {
                            if (typeof cb === "function") cb.apply(this, arguments);
                            retry();
                        }
                    }));
                }
            }
            return fs;
        }
        function enqueue(elem) {
            debug("ENQUEUE", elem[0].name, elem[1]);
            fs[gracefulQueue].push(elem);
        }
        function retry() {
            var elem = fs[gracefulQueue].shift();
            if (elem) {
                debug("RETRY", elem[0].name, elem[1]);
                elem[0].apply(null, elem[1]);
            }
        }
    },
    7885: (module, __unused_webpack_exports, __webpack_require__) => {
        var Stream = __webpack_require__(2413).Stream;
        module.exports = legacy;
        function legacy(fs) {
            return {
                ReadStream,
                WriteStream
            };
            function ReadStream(path, options) {
                if (!(this instanceof ReadStream)) return new ReadStream(path, options);
                Stream.call(this);
                var self = this;
                this.path = path;
                this.fd = null;
                this.readable = true;
                this.paused = false;
                this.flags = "r";
                this.mode = 438;
                this.bufferSize = 64 * 1024;
                options = options || {};
                var keys = Object.keys(options);
                for (var index = 0, length = keys.length; index < length; index++) {
                    var key = keys[index];
                    this[key] = options[key];
                }
                if (this.encoding) this.setEncoding(this.encoding);
                if (this.start !== undefined) {
                    if ("number" !== typeof this.start) {
                        throw TypeError("start must be a Number");
                    }
                    if (this.end === undefined) {
                        this.end = Infinity;
                    } else if ("number" !== typeof this.end) {
                        throw TypeError("end must be a Number");
                    }
                    if (this.start > this.end) {
                        throw new Error("start must be <= end");
                    }
                    this.pos = this.start;
                }
                if (this.fd !== null) {
                    process.nextTick((function() {
                        self._read();
                    }));
                    return;
                }
                fs.open(this.path, this.flags, this.mode, (function(err, fd) {
                    if (err) {
                        self.emit("error", err);
                        self.readable = false;
                        return;
                    }
                    self.fd = fd;
                    self.emit("open", fd);
                    self._read();
                }));
            }
            function WriteStream(path, options) {
                if (!(this instanceof WriteStream)) return new WriteStream(path, options);
                Stream.call(this);
                this.path = path;
                this.fd = null;
                this.writable = true;
                this.flags = "w";
                this.encoding = "binary";
                this.mode = 438;
                this.bytesWritten = 0;
                options = options || {};
                var keys = Object.keys(options);
                for (var index = 0, length = keys.length; index < length; index++) {
                    var key = keys[index];
                    this[key] = options[key];
                }
                if (this.start !== undefined) {
                    if ("number" !== typeof this.start) {
                        throw TypeError("start must be a Number");
                    }
                    if (this.start < 0) {
                        throw new Error("start must be >= zero");
                    }
                    this.pos = this.start;
                }
                this.busy = false;
                this._queue = [];
                if (this.fd === null) {
                    this._open = fs.open;
                    this._queue.push([ this._open, this.path, this.flags, this.mode, undefined ]);
                    this.flush();
                }
            }
        }
    },
    7994: (module, __unused_webpack_exports, __webpack_require__) => {
        var constants = __webpack_require__(7619);
        var origCwd = process.cwd;
        var cwd = null;
        var platform = process.env.GRACEFUL_FS_PLATFORM || process.platform;
        process.cwd = function() {
            if (!cwd) cwd = origCwd.call(process);
            return cwd;
        };
        try {
            process.cwd();
        } catch (er) {}
        if (typeof process.chdir === "function") {
            var chdir = process.chdir;
            process.chdir = function(d) {
                cwd = null;
                chdir.call(process, d);
            };
            if (Object.setPrototypeOf) Object.setPrototypeOf(process.chdir, chdir);
        }
        module.exports = patch;
        function patch(fs) {
            if (constants.hasOwnProperty("O_SYMLINK") && process.version.match(/^v0\.6\.[0-2]|^v0\.5\./)) {
                patchLchmod(fs);
            }
            if (!fs.lutimes) {
                patchLutimes(fs);
            }
            fs.chown = chownFix(fs.chown);
            fs.fchown = chownFix(fs.fchown);
            fs.lchown = chownFix(fs.lchown);
            fs.chmod = chmodFix(fs.chmod);
            fs.fchmod = chmodFix(fs.fchmod);
            fs.lchmod = chmodFix(fs.lchmod);
            fs.chownSync = chownFixSync(fs.chownSync);
            fs.fchownSync = chownFixSync(fs.fchownSync);
            fs.lchownSync = chownFixSync(fs.lchownSync);
            fs.chmodSync = chmodFixSync(fs.chmodSync);
            fs.fchmodSync = chmodFixSync(fs.fchmodSync);
            fs.lchmodSync = chmodFixSync(fs.lchmodSync);
            fs.stat = statFix(fs.stat);
            fs.fstat = statFix(fs.fstat);
            fs.lstat = statFix(fs.lstat);
            fs.statSync = statFixSync(fs.statSync);
            fs.fstatSync = statFixSync(fs.fstatSync);
            fs.lstatSync = statFixSync(fs.lstatSync);
            if (!fs.lchmod) {
                fs.lchmod = function(path, mode, cb) {
                    if (cb) process.nextTick(cb);
                };
                fs.lchmodSync = function() {};
            }
            if (!fs.lchown) {
                fs.lchown = function(path, uid, gid, cb) {
                    if (cb) process.nextTick(cb);
                };
                fs.lchownSync = function() {};
            }
            if (platform === "win32") {
                fs.rename = function(fs$rename) {
                    return function(from, to, cb) {
                        var start = Date.now();
                        var backoff = 0;
                        fs$rename(from, to, (function CB(er) {
                            if (er && (er.code === "EACCES" || er.code === "EPERM") && Date.now() - start < 6e4) {
                                setTimeout((function() {
                                    fs.stat(to, (function(stater, st) {
                                        if (stater && stater.code === "ENOENT") fs$rename(from, to, CB); else cb(er);
                                    }));
                                }), backoff);
                                if (backoff < 100) backoff += 10;
                                return;
                            }
                            if (cb) cb(er);
                        }));
                    };
                }(fs.rename);
            }
            fs.read = function(fs$read) {
                function read(fd, buffer, offset, length, position, callback_) {
                    var callback;
                    if (callback_ && typeof callback_ === "function") {
                        var eagCounter = 0;
                        callback = function(er, _, __) {
                            if (er && er.code === "EAGAIN" && eagCounter < 10) {
                                eagCounter++;
                                return fs$read.call(fs, fd, buffer, offset, length, position, callback);
                            }
                            callback_.apply(this, arguments);
                        };
                    }
                    return fs$read.call(fs, fd, buffer, offset, length, position, callback);
                }
                if (Object.setPrototypeOf) Object.setPrototypeOf(read, fs$read);
                return read;
            }(fs.read);
            fs.readSync = function(fs$readSync) {
                return function(fd, buffer, offset, length, position) {
                    var eagCounter = 0;
                    while (true) {
                        try {
                            return fs$readSync.call(fs, fd, buffer, offset, length, position);
                        } catch (er) {
                            if (er.code === "EAGAIN" && eagCounter < 10) {
                                eagCounter++;
                                continue;
                            }
                            throw er;
                        }
                    }
                };
            }(fs.readSync);
            function patchLchmod(fs) {
                fs.lchmod = function(path, mode, callback) {
                    fs.open(path, constants.O_WRONLY | constants.O_SYMLINK, mode, (function(err, fd) {
                        if (err) {
                            if (callback) callback(err);
                            return;
                        }
                        fs.fchmod(fd, mode, (function(err) {
                            fs.close(fd, (function(err2) {
                                if (callback) callback(err || err2);
                            }));
                        }));
                    }));
                };
                fs.lchmodSync = function(path, mode) {
                    var fd = fs.openSync(path, constants.O_WRONLY | constants.O_SYMLINK, mode);
                    var threw = true;
                    var ret;
                    try {
                        ret = fs.fchmodSync(fd, mode);
                        threw = false;
                    } finally {
                        if (threw) {
                            try {
                                fs.closeSync(fd);
                            } catch (er) {}
                        } else {
                            fs.closeSync(fd);
                        }
                    }
                    return ret;
                };
            }
            function patchLutimes(fs) {
                if (constants.hasOwnProperty("O_SYMLINK")) {
                    fs.lutimes = function(path, at, mt, cb) {
                        fs.open(path, constants.O_SYMLINK, (function(er, fd) {
                            if (er) {
                                if (cb) cb(er);
                                return;
                            }
                            fs.futimes(fd, at, mt, (function(er) {
                                fs.close(fd, (function(er2) {
                                    if (cb) cb(er || er2);
                                }));
                            }));
                        }));
                    };
                    fs.lutimesSync = function(path, at, mt) {
                        var fd = fs.openSync(path, constants.O_SYMLINK);
                        var ret;
                        var threw = true;
                        try {
                            ret = fs.futimesSync(fd, at, mt);
                            threw = false;
                        } finally {
                            if (threw) {
                                try {
                                    fs.closeSync(fd);
                                } catch (er) {}
                            } else {
                                fs.closeSync(fd);
                            }
                        }
                        return ret;
                    };
                } else {
                    fs.lutimes = function(_a, _b, _c, cb) {
                        if (cb) process.nextTick(cb);
                    };
                    fs.lutimesSync = function() {};
                }
            }
            function chmodFix(orig) {
                if (!orig) return orig;
                return function(target, mode, cb) {
                    return orig.call(fs, target, mode, (function(er) {
                        if (chownErOk(er)) er = null;
                        if (cb) cb.apply(this, arguments);
                    }));
                };
            }
            function chmodFixSync(orig) {
                if (!orig) return orig;
                return function(target, mode) {
                    try {
                        return orig.call(fs, target, mode);
                    } catch (er) {
                        if (!chownErOk(er)) throw er;
                    }
                };
            }
            function chownFix(orig) {
                if (!orig) return orig;
                return function(target, uid, gid, cb) {
                    return orig.call(fs, target, uid, gid, (function(er) {
                        if (chownErOk(er)) er = null;
                        if (cb) cb.apply(this, arguments);
                    }));
                };
            }
            function chownFixSync(orig) {
                if (!orig) return orig;
                return function(target, uid, gid) {
                    try {
                        return orig.call(fs, target, uid, gid);
                    } catch (er) {
                        if (!chownErOk(er)) throw er;
                    }
                };
            }
            function statFix(orig) {
                if (!orig) return orig;
                return function(target, options, cb) {
                    if (typeof options === "function") {
                        cb = options;
                        options = null;
                    }
                    function callback(er, stats) {
                        if (stats) {
                            if (stats.uid < 0) stats.uid += 4294967296;
                            if (stats.gid < 0) stats.gid += 4294967296;
                        }
                        if (cb) cb.apply(this, arguments);
                    }
                    return options ? orig.call(fs, target, options, callback) : orig.call(fs, target, callback);
                };
            }
            function statFixSync(orig) {
                if (!orig) return orig;
                return function(target, options) {
                    var stats = options ? orig.call(fs, target, options) : orig.call(fs, target);
                    if (stats.uid < 0) stats.uid += 4294967296;
                    if (stats.gid < 0) stats.gid += 4294967296;
                    return stats;
                };
            }
            function chownErOk(er) {
                if (!er) return true;
                if (er.code === "ENOSYS") return true;
                var nonroot = !process.getuid || process.getuid() !== 0;
                if (nonroot) {
                    if (er.code === "EINVAL" || er.code === "EPERM") return true;
                }
                return false;
            }
        }
    },
    3393: (module, __unused_webpack_exports, __webpack_require__) => {
        let _fs;
        try {
            _fs = __webpack_require__(6851);
        } catch (_) {
            _fs = __webpack_require__(5747);
        }
        const universalify = __webpack_require__(6267);
        const {stringify, stripBom} = __webpack_require__(9293);
        async function _readFile(file, options = {}) {
            if (typeof options === "string") {
                options = {
                    encoding: options
                };
            }
            const fs = options.fs || _fs;
            const shouldThrow = "throws" in options ? options.throws : true;
            let data = await universalify.fromCallback(fs.readFile)(file, options);
            data = stripBom(data);
            let obj;
            try {
                obj = JSON.parse(data, options ? options.reviver : null);
            } catch (err) {
                if (shouldThrow) {
                    err.message = `${file}: ${err.message}`;
                    throw err;
                } else {
                    return null;
                }
            }
            return obj;
        }
        const readFile = universalify.fromPromise(_readFile);
        function readFileSync(file, options = {}) {
            if (typeof options === "string") {
                options = {
                    encoding: options
                };
            }
            const fs = options.fs || _fs;
            const shouldThrow = "throws" in options ? options.throws : true;
            try {
                let content = fs.readFileSync(file, options);
                content = stripBom(content);
                return JSON.parse(content, options.reviver);
            } catch (err) {
                if (shouldThrow) {
                    err.message = `${file}: ${err.message}`;
                    throw err;
                } else {
                    return null;
                }
            }
        }
        async function _writeFile(file, obj, options = {}) {
            const fs = options.fs || _fs;
            const str = stringify(obj, options);
            await universalify.fromCallback(fs.writeFile)(file, str, options);
        }
        const writeFile = universalify.fromPromise(_writeFile);
        function writeFileSync(file, obj, options = {}) {
            const fs = options.fs || _fs;
            const str = stringify(obj, options);
            return fs.writeFileSync(file, str, options);
        }
        const jsonfile = {
            readFile,
            readFileSync,
            writeFile,
            writeFileSync
        };
        module.exports = jsonfile;
    },
    6267: (__unused_webpack_module, exports) => {
        "use strict";
        exports.fromCallback = function(fn) {
            return Object.defineProperty((function(...args) {
                if (typeof args[args.length - 1] === "function") fn.apply(this, args); else {
                    return new Promise(((resolve, reject) => {
                        fn.call(this, ...args, ((err, res) => err != null ? reject(err) : resolve(res)));
                    }));
                }
            }), "name", {
                value: fn.name
            });
        };
        exports.fromPromise = function(fn) {
            return Object.defineProperty((function(...args) {
                const cb = args[args.length - 1];
                if (typeof cb !== "function") return fn.apply(this, args); else fn.apply(this, args.slice(0, -1)).then((r => cb(null, r)), cb);
            }), "name", {
                value: fn.name
            });
        };
    },
    9293: module => {
        function stringify(obj, {EOL = "\n", finalEOL = true, replacer = null, spaces} = {}) {
            const EOF = finalEOL ? EOL : "";
            const str = JSON.stringify(obj, replacer, spaces);
            return str.replace(/\n/g, EOL) + EOF;
        }
        function stripBom(content) {
            if (Buffer.isBuffer(content)) content = content.toString("utf8");
            return content.replace(/^\uFEFF/, "");
        }
        module.exports = {
            stringify,
            stripBom
        };
    },
    6120: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        var helpers = __webpack_require__(6495);
        var ValidatorResult = helpers.ValidatorResult;
        var SchemaError = helpers.SchemaError;
        var attribute = {};
        attribute.ignoreProperties = {
            id: true,
            default: true,
            description: true,
            title: true,
            additionalItems: true,
            then: true,
            else: true,
            $schema: true,
            $ref: true,
            extends: true
        };
        var validators = attribute.validators = {};
        validators.type = function validateType(instance, schema, options, ctx) {
            if (instance === undefined) {
                return null;
            }
            var result = new ValidatorResult(instance, schema, options, ctx);
            var types = Array.isArray(schema.type) ? schema.type : [ schema.type ];
            if (!types.some(this.testType.bind(this, instance, schema, options, ctx))) {
                var list = types.map((function(v) {
                    if (!v) return;
                    var id = v.$id || v.id;
                    return id ? "<" + id + ">" : v + "";
                }));
                result.addError({
                    name: "type",
                    argument: list,
                    message: "is not of a type(s) " + list
                });
            }
            return result;
        };
        function testSchemaNoThrow(instance, options, ctx, callback, schema) {
            var throwError = options.throwError;
            var throwAll = options.throwAll;
            options.throwError = false;
            options.throwAll = false;
            var res = this.validateSchema(instance, schema, options, ctx);
            options.throwError = throwError;
            options.throwAll = throwAll;
            if (!res.valid && callback instanceof Function) {
                callback(res);
            }
            return res.valid;
        }
        validators.anyOf = function validateAnyOf(instance, schema, options, ctx) {
            if (instance === undefined) {
                return null;
            }
            var result = new ValidatorResult(instance, schema, options, ctx);
            var inner = new ValidatorResult(instance, schema, options, ctx);
            if (!Array.isArray(schema.anyOf)) {
                throw new SchemaError("anyOf must be an array");
            }
            if (!schema.anyOf.some(testSchemaNoThrow.bind(this, instance, options, ctx, (function(res) {
                inner.importErrors(res);
            })))) {
                var list = schema.anyOf.map((function(v, i) {
                    var id = v.$id || v.id;
                    if (id) return "<" + id + ">";
                    return v.title && JSON.stringify(v.title) || v["$ref"] && "<" + v["$ref"] + ">" || "[subschema " + i + "]";
                }));
                if (options.nestedErrors) {
                    result.importErrors(inner);
                }
                result.addError({
                    name: "anyOf",
                    argument: list,
                    message: "is not any of " + list.join(",")
                });
            }
            return result;
        };
        validators.allOf = function validateAllOf(instance, schema, options, ctx) {
            if (instance === undefined) {
                return null;
            }
            if (!Array.isArray(schema.allOf)) {
                throw new SchemaError("allOf must be an array");
            }
            var result = new ValidatorResult(instance, schema, options, ctx);
            var self = this;
            schema.allOf.forEach((function(v, i) {
                var valid = self.validateSchema(instance, v, options, ctx);
                if (!valid.valid) {
                    var id = v.$id || v.id;
                    var msg = id || v.title && JSON.stringify(v.title) || v["$ref"] && "<" + v["$ref"] + ">" || "[subschema " + i + "]";
                    result.addError({
                        name: "allOf",
                        argument: {
                            id: msg,
                            length: valid.errors.length,
                            valid
                        },
                        message: "does not match allOf schema " + msg + " with " + valid.errors.length + " error[s]:"
                    });
                    result.importErrors(valid);
                }
            }));
            return result;
        };
        validators.oneOf = function validateOneOf(instance, schema, options, ctx) {
            if (instance === undefined) {
                return null;
            }
            if (!Array.isArray(schema.oneOf)) {
                throw new SchemaError("oneOf must be an array");
            }
            var result = new ValidatorResult(instance, schema, options, ctx);
            var inner = new ValidatorResult(instance, schema, options, ctx);
            var count = schema.oneOf.filter(testSchemaNoThrow.bind(this, instance, options, ctx, (function(res) {
                inner.importErrors(res);
            }))).length;
            var list = schema.oneOf.map((function(v, i) {
                var id = v.$id || v.id;
                return id || v.title && JSON.stringify(v.title) || v["$ref"] && "<" + v["$ref"] + ">" || "[subschema " + i + "]";
            }));
            if (count !== 1) {
                if (options.nestedErrors) {
                    result.importErrors(inner);
                }
                result.addError({
                    name: "oneOf",
                    argument: list,
                    message: "is not exactly one from " + list.join(",")
                });
            }
            return result;
        };
        validators.if = function validateIf(instance, schema, options, ctx) {
            if (instance === undefined) return null;
            if (!helpers.isSchema(schema.if)) throw new Error('Expected "if" keyword to be a schema');
            var ifValid = testSchemaNoThrow.call(this, instance, options, ctx, null, schema.if);
            var result = new ValidatorResult(instance, schema, options, ctx);
            var res;
            if (ifValid) {
                if (schema.then === undefined) return;
                if (!helpers.isSchema(schema.then)) throw new Error('Expected "then" keyword to be a schema');
                res = this.validateSchema(instance, schema.then, options, ctx.makeChild(schema.then));
                result.importErrors(res);
            } else {
                if (schema.else === undefined) return;
                if (!helpers.isSchema(schema.else)) throw new Error('Expected "else" keyword to be a schema');
                res = this.validateSchema(instance, schema.else, options, ctx.makeChild(schema.else));
                result.importErrors(res);
            }
            return result;
        };
        function getEnumerableProperty(object, key) {
            if (Object.hasOwnProperty.call(object, key)) return object[key];
            if (!(key in object)) return;
            while (object = Object.getPrototypeOf(object)) {
                if (Object.propertyIsEnumerable.call(object, key)) return object[key];
            }
        }
        validators.propertyNames = function validatePropertyNames(instance, schema, options, ctx) {
            if (!this.types.object(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var subschema = schema.propertyNames !== undefined ? schema.propertyNames : {};
            if (!helpers.isSchema(subschema)) throw new SchemaError('Expected "propertyNames" to be a schema (object or boolean)');
            for (var property in instance) {
                if (getEnumerableProperty(instance, property) !== undefined) {
                    var res = this.validateSchema(property, subschema, options, ctx.makeChild(subschema));
                    result.importErrors(res);
                }
            }
            return result;
        };
        validators.properties = function validateProperties(instance, schema, options, ctx) {
            if (!this.types.object(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var properties = schema.properties || {};
            for (var property in properties) {
                var subschema = properties[property];
                if (subschema === undefined) {
                    continue;
                } else if (subschema === null) {
                    throw new SchemaError('Unexpected null, expected schema in "properties"');
                }
                if (typeof options.preValidateProperty == "function") {
                    options.preValidateProperty(instance, property, subschema, options, ctx);
                }
                var prop = getEnumerableProperty(instance, property);
                var res = this.validateSchema(prop, subschema, options, ctx.makeChild(subschema, property));
                if (res.instance !== result.instance[property]) result.instance[property] = res.instance;
                result.importErrors(res);
            }
            return result;
        };
        function testAdditionalProperty(instance, schema, options, ctx, property, result) {
            if (!this.types.object(instance)) return;
            if (schema.properties && schema.properties[property] !== undefined) {
                return;
            }
            if (schema.additionalProperties === false) {
                result.addError({
                    name: "additionalProperties",
                    argument: property,
                    message: "is not allowed to have the additional property " + JSON.stringify(property)
                });
            } else {
                var additionalProperties = schema.additionalProperties || {};
                if (typeof options.preValidateProperty == "function") {
                    options.preValidateProperty(instance, property, additionalProperties, options, ctx);
                }
                var res = this.validateSchema(instance[property], additionalProperties, options, ctx.makeChild(additionalProperties, property));
                if (res.instance !== result.instance[property]) result.instance[property] = res.instance;
                result.importErrors(res);
            }
        }
        validators.patternProperties = function validatePatternProperties(instance, schema, options, ctx) {
            if (!this.types.object(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var patternProperties = schema.patternProperties || {};
            for (var property in instance) {
                var test = true;
                for (var pattern in patternProperties) {
                    var subschema = patternProperties[pattern];
                    if (subschema === undefined) {
                        continue;
                    } else if (subschema === null) {
                        throw new SchemaError('Unexpected null, expected schema in "patternProperties"');
                    }
                    try {
                        var regexp = new RegExp(pattern, "u");
                    } catch (_e) {
                        regexp = new RegExp(pattern);
                    }
                    if (!regexp.test(property)) {
                        continue;
                    }
                    test = false;
                    if (typeof options.preValidateProperty == "function") {
                        options.preValidateProperty(instance, property, subschema, options, ctx);
                    }
                    var res = this.validateSchema(instance[property], subschema, options, ctx.makeChild(subschema, property));
                    if (res.instance !== result.instance[property]) result.instance[property] = res.instance;
                    result.importErrors(res);
                }
                if (test) {
                    testAdditionalProperty.call(this, instance, schema, options, ctx, property, result);
                }
            }
            return result;
        };
        validators.additionalProperties = function validateAdditionalProperties(instance, schema, options, ctx) {
            if (!this.types.object(instance)) return;
            if (schema.patternProperties) {
                return null;
            }
            var result = new ValidatorResult(instance, schema, options, ctx);
            for (var property in instance) {
                testAdditionalProperty.call(this, instance, schema, options, ctx, property, result);
            }
            return result;
        };
        validators.minProperties = function validateMinProperties(instance, schema, options, ctx) {
            if (!this.types.object(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var keys = Object.keys(instance);
            if (!(keys.length >= schema.minProperties)) {
                result.addError({
                    name: "minProperties",
                    argument: schema.minProperties,
                    message: "does not meet minimum property length of " + schema.minProperties
                });
            }
            return result;
        };
        validators.maxProperties = function validateMaxProperties(instance, schema, options, ctx) {
            if (!this.types.object(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var keys = Object.keys(instance);
            if (!(keys.length <= schema.maxProperties)) {
                result.addError({
                    name: "maxProperties",
                    argument: schema.maxProperties,
                    message: "does not meet maximum property length of " + schema.maxProperties
                });
            }
            return result;
        };
        validators.items = function validateItems(instance, schema, options, ctx) {
            var self = this;
            if (!this.types.array(instance)) return;
            if (!schema.items) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            instance.every((function(value, i) {
                var items = Array.isArray(schema.items) ? schema.items[i] || schema.additionalItems : schema.items;
                if (items === undefined) {
                    return true;
                }
                if (items === false) {
                    result.addError({
                        name: "items",
                        message: "additionalItems not permitted"
                    });
                    return false;
                }
                var res = self.validateSchema(value, items, options, ctx.makeChild(items, i));
                if (res.instance !== result.instance[i]) result.instance[i] = res.instance;
                result.importErrors(res);
                return true;
            }));
            return result;
        };
        validators.minimum = function validateMinimum(instance, schema, options, ctx) {
            if (!this.types.number(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (schema.exclusiveMinimum && schema.exclusiveMinimum === true) {
                if (!(instance > schema.minimum)) {
                    result.addError({
                        name: "minimum",
                        argument: schema.minimum,
                        message: "must be greater than " + schema.minimum
                    });
                }
            } else {
                if (!(instance >= schema.minimum)) {
                    result.addError({
                        name: "minimum",
                        argument: schema.minimum,
                        message: "must be greater than or equal to " + schema.minimum
                    });
                }
            }
            return result;
        };
        validators.maximum = function validateMaximum(instance, schema, options, ctx) {
            if (!this.types.number(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (schema.exclusiveMaximum && schema.exclusiveMaximum === true) {
                if (!(instance < schema.maximum)) {
                    result.addError({
                        name: "maximum",
                        argument: schema.maximum,
                        message: "must be less than " + schema.maximum
                    });
                }
            } else {
                if (!(instance <= schema.maximum)) {
                    result.addError({
                        name: "maximum",
                        argument: schema.maximum,
                        message: "must be less than or equal to " + schema.maximum
                    });
                }
            }
            return result;
        };
        validators.exclusiveMinimum = function validateExclusiveMinimum(instance, schema, options, ctx) {
            if (typeof schema.exclusiveMaximum === "boolean") return;
            if (!this.types.number(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var valid = instance > schema.exclusiveMinimum;
            if (!valid) {
                result.addError({
                    name: "exclusiveMinimum",
                    argument: schema.exclusiveMinimum,
                    message: "must be strictly greater than " + schema.exclusiveMinimum
                });
            }
            return result;
        };
        validators.exclusiveMaximum = function validateExclusiveMaximum(instance, schema, options, ctx) {
            if (typeof schema.exclusiveMaximum === "boolean") return;
            if (!this.types.number(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var valid = instance < schema.exclusiveMaximum;
            if (!valid) {
                result.addError({
                    name: "exclusiveMaximum",
                    argument: schema.exclusiveMaximum,
                    message: "must be strictly less than " + schema.exclusiveMaximum
                });
            }
            return result;
        };
        var validateMultipleOfOrDivisbleBy = function validateMultipleOfOrDivisbleBy(instance, schema, options, ctx, validationType, errorMessage) {
            if (!this.types.number(instance)) return;
            var validationArgument = schema[validationType];
            if (validationArgument == 0) {
                throw new SchemaError(validationType + " cannot be zero");
            }
            var result = new ValidatorResult(instance, schema, options, ctx);
            var instanceDecimals = helpers.getDecimalPlaces(instance);
            var divisorDecimals = helpers.getDecimalPlaces(validationArgument);
            var maxDecimals = Math.max(instanceDecimals, divisorDecimals);
            var multiplier = Math.pow(10, maxDecimals);
            if (Math.round(instance * multiplier) % Math.round(validationArgument * multiplier) !== 0) {
                result.addError({
                    name: validationType,
                    argument: validationArgument,
                    message: errorMessage + JSON.stringify(validationArgument)
                });
            }
            return result;
        };
        validators.multipleOf = function validateMultipleOf(instance, schema, options, ctx) {
            return validateMultipleOfOrDivisbleBy.call(this, instance, schema, options, ctx, "multipleOf", "is not a multiple of (divisible by) ");
        };
        validators.divisibleBy = function validateDivisibleBy(instance, schema, options, ctx) {
            return validateMultipleOfOrDivisbleBy.call(this, instance, schema, options, ctx, "divisibleBy", "is not divisible by (multiple of) ");
        };
        validators.required = function validateRequired(instance, schema, options, ctx) {
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (instance === undefined && schema.required === true) {
                result.addError({
                    name: "required",
                    message: "is required"
                });
            } else if (this.types.object(instance) && Array.isArray(schema.required)) {
                schema.required.forEach((function(n) {
                    if (getEnumerableProperty(instance, n) === undefined) {
                        result.addError({
                            name: "required",
                            argument: n,
                            message: "requires property " + JSON.stringify(n)
                        });
                    }
                }));
            }
            return result;
        };
        validators.pattern = function validatePattern(instance, schema, options, ctx) {
            if (!this.types.string(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var pattern = schema.pattern;
            try {
                var regexp = new RegExp(pattern, "u");
            } catch (_e) {
                regexp = new RegExp(pattern);
            }
            if (!instance.match(regexp)) {
                result.addError({
                    name: "pattern",
                    argument: schema.pattern,
                    message: "does not match pattern " + JSON.stringify(schema.pattern.toString())
                });
            }
            return result;
        };
        validators.format = function validateFormat(instance, schema, options, ctx) {
            if (instance === undefined) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (!result.disableFormat && !helpers.isFormat(instance, schema.format, this)) {
                result.addError({
                    name: "format",
                    argument: schema.format,
                    message: "does not conform to the " + JSON.stringify(schema.format) + " format"
                });
            }
            return result;
        };
        validators.minLength = function validateMinLength(instance, schema, options, ctx) {
            if (!this.types.string(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var hsp = instance.match(/[\uDC00-\uDFFF]/g);
            var length = instance.length - (hsp ? hsp.length : 0);
            if (!(length >= schema.minLength)) {
                result.addError({
                    name: "minLength",
                    argument: schema.minLength,
                    message: "does not meet minimum length of " + schema.minLength
                });
            }
            return result;
        };
        validators.maxLength = function validateMaxLength(instance, schema, options, ctx) {
            if (!this.types.string(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var hsp = instance.match(/[\uDC00-\uDFFF]/g);
            var length = instance.length - (hsp ? hsp.length : 0);
            if (!(length <= schema.maxLength)) {
                result.addError({
                    name: "maxLength",
                    argument: schema.maxLength,
                    message: "does not meet maximum length of " + schema.maxLength
                });
            }
            return result;
        };
        validators.minItems = function validateMinItems(instance, schema, options, ctx) {
            if (!this.types.array(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (!(instance.length >= schema.minItems)) {
                result.addError({
                    name: "minItems",
                    argument: schema.minItems,
                    message: "does not meet minimum length of " + schema.minItems
                });
            }
            return result;
        };
        validators.maxItems = function validateMaxItems(instance, schema, options, ctx) {
            if (!this.types.array(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (!(instance.length <= schema.maxItems)) {
                result.addError({
                    name: "maxItems",
                    argument: schema.maxItems,
                    message: "does not meet maximum length of " + schema.maxItems
                });
            }
            return result;
        };
        function testArrays(v, i, a) {
            var j, len = a.length;
            for (j = i + 1, len; j < len; j++) {
                if (helpers.deepCompareStrict(v, a[j])) {
                    return false;
                }
            }
            return true;
        }
        validators.uniqueItems = function validateUniqueItems(instance, schema, options, ctx) {
            if (schema.uniqueItems !== true) return;
            if (!this.types.array(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (!instance.every(testArrays)) {
                result.addError({
                    name: "uniqueItems",
                    message: "contains duplicate item"
                });
            }
            return result;
        };
        validators.dependencies = function validateDependencies(instance, schema, options, ctx) {
            if (!this.types.object(instance)) return;
            var result = new ValidatorResult(instance, schema, options, ctx);
            for (var property in schema.dependencies) {
                if (instance[property] === undefined) {
                    continue;
                }
                var dep = schema.dependencies[property];
                var childContext = ctx.makeChild(dep, property);
                if (typeof dep == "string") {
                    dep = [ dep ];
                }
                if (Array.isArray(dep)) {
                    dep.forEach((function(prop) {
                        if (instance[prop] === undefined) {
                            result.addError({
                                name: "dependencies",
                                argument: childContext.propertyPath,
                                message: "property " + prop + " not found, required by " + childContext.propertyPath
                            });
                        }
                    }));
                } else {
                    var res = this.validateSchema(instance, dep, options, childContext);
                    if (result.instance !== res.instance) result.instance = res.instance;
                    if (res && res.errors.length) {
                        result.addError({
                            name: "dependencies",
                            argument: childContext.propertyPath,
                            message: "does not meet dependency required by " + childContext.propertyPath
                        });
                        result.importErrors(res);
                    }
                }
            }
            return result;
        };
        validators["enum"] = function validateEnum(instance, schema, options, ctx) {
            if (instance === undefined) {
                return null;
            }
            if (!Array.isArray(schema["enum"])) {
                throw new SchemaError("enum expects an array", schema);
            }
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (!schema["enum"].some(helpers.deepCompareStrict.bind(null, instance))) {
                result.addError({
                    name: "enum",
                    argument: schema["enum"],
                    message: "is not one of enum values: " + schema["enum"].map(String).join(",")
                });
            }
            return result;
        };
        validators["const"] = function validateEnum(instance, schema, options, ctx) {
            if (instance === undefined) {
                return null;
            }
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (!helpers.deepCompareStrict(schema["const"], instance)) {
                result.addError({
                    name: "const",
                    argument: schema["const"],
                    message: "does not exactly match expected constant: " + schema["const"]
                });
            }
            return result;
        };
        validators.not = validators.disallow = function validateNot(instance, schema, options, ctx) {
            var self = this;
            if (instance === undefined) return null;
            var result = new ValidatorResult(instance, schema, options, ctx);
            var notTypes = schema.not || schema.disallow;
            if (!notTypes) return null;
            if (!Array.isArray(notTypes)) notTypes = [ notTypes ];
            notTypes.forEach((function(type) {
                if (self.testType(instance, schema, options, ctx, type)) {
                    var id = type && (type.$id || type.id);
                    var schemaId = id || type;
                    result.addError({
                        name: "not",
                        argument: schemaId,
                        message: "is of prohibited type " + schemaId
                    });
                }
            }));
            return result;
        };
        module.exports = attribute;
    },
    6495: (module, exports, __webpack_require__) => {
        "use strict";
        var uri = __webpack_require__(8835);
        var ValidationError = exports.ValidationError = function ValidationError(message, instance, schema, path, name, argument) {
            if (Array.isArray(path)) {
                this.path = path;
                this.property = path.reduce((function(sum, item) {
                    return sum + makeSuffix(item);
                }), "instance");
            } else if (path !== undefined) {
                this.property = path;
            }
            if (message) {
                this.message = message;
            }
            if (schema) {
                var id = schema.$id || schema.id;
                this.schema = id || schema;
            }
            if (instance !== undefined) {
                this.instance = instance;
            }
            this.name = name;
            this.argument = argument;
            this.stack = this.toString();
        };
        ValidationError.prototype.toString = function toString() {
            return this.property + " " + this.message;
        };
        var ValidatorResult = exports.ValidatorResult = function ValidatorResult(instance, schema, options, ctx) {
            this.instance = instance;
            this.schema = schema;
            this.options = options;
            this.path = ctx.path;
            this.propertyPath = ctx.propertyPath;
            this.errors = [];
            this.throwError = options && options.throwError;
            this.throwFirst = options && options.throwFirst;
            this.throwAll = options && options.throwAll;
            this.disableFormat = options && options.disableFormat === true;
        };
        ValidatorResult.prototype.addError = function addError(detail) {
            var err;
            if (typeof detail == "string") {
                err = new ValidationError(detail, this.instance, this.schema, this.path);
            } else {
                if (!detail) throw new Error("Missing error detail");
                if (!detail.message) throw new Error("Missing error message");
                if (!detail.name) throw new Error("Missing validator type");
                err = new ValidationError(detail.message, this.instance, this.schema, this.path, detail.name, detail.argument);
            }
            this.errors.push(err);
            if (this.throwFirst) {
                throw new ValidatorResultError(this);
            } else if (this.throwError) {
                throw err;
            }
            return err;
        };
        ValidatorResult.prototype.importErrors = function importErrors(res) {
            if (typeof res == "string" || res && res.validatorType) {
                this.addError(res);
            } else if (res && res.errors) {
                Array.prototype.push.apply(this.errors, res.errors);
            }
        };
        function stringizer(v, i) {
            return i + ": " + v.toString() + "\n";
        }
        ValidatorResult.prototype.toString = function toString(res) {
            return this.errors.map(stringizer).join("");
        };
        Object.defineProperty(ValidatorResult.prototype, "valid", {
            get: function() {
                return !this.errors.length;
            }
        });
        module.exports.ValidatorResultError = ValidatorResultError;
        function ValidatorResultError(result) {
            if (Error.captureStackTrace) {
                Error.captureStackTrace(this, ValidatorResultError);
            }
            this.instance = result.instance;
            this.schema = result.schema;
            this.options = result.options;
            this.errors = result.errors;
        }
        ValidatorResultError.prototype = new Error;
        ValidatorResultError.prototype.constructor = ValidatorResultError;
        ValidatorResultError.prototype.name = "Validation Error";
        var SchemaError = exports.SchemaError = function SchemaError(msg, schema) {
            this.message = msg;
            this.schema = schema;
            Error.call(this, msg);
            Error.captureStackTrace(this, SchemaError);
        };
        SchemaError.prototype = Object.create(Error.prototype, {
            constructor: {
                value: SchemaError,
                enumerable: false
            },
            name: {
                value: "SchemaError",
                enumerable: false
            }
        });
        var SchemaContext = exports.SchemaContext = function SchemaContext(schema, options, path, base, schemas) {
            this.schema = schema;
            this.options = options;
            if (Array.isArray(path)) {
                this.path = path;
                this.propertyPath = path.reduce((function(sum, item) {
                    return sum + makeSuffix(item);
                }), "instance");
            } else {
                this.propertyPath = path;
            }
            this.base = base;
            this.schemas = schemas;
        };
        SchemaContext.prototype.resolve = function resolve(target) {
            return uri.resolve(this.base, target);
        };
        SchemaContext.prototype.makeChild = function makeChild(schema, propertyName) {
            var path = propertyName === undefined ? this.path : this.path.concat([ propertyName ]);
            var id = schema.$id || schema.id;
            var base = uri.resolve(this.base, id || "");
            var ctx = new SchemaContext(schema, this.options, path, base, Object.create(this.schemas));
            if (id && !ctx.schemas[base]) {
                ctx.schemas[base] = schema;
            }
            return ctx;
        };
        var FORMAT_REGEXPS = exports.FORMAT_REGEXPS = {
            "date-time": /^\d{4}-(?:0[0-9]{1}|1[0-2]{1})-(3[01]|0[1-9]|[12][0-9])[tT ](2[0-4]|[01][0-9]):([0-5][0-9]):(60|[0-5][0-9])(\.\d+)?([zZ]|[+-]([0-5][0-9]):(60|[0-5][0-9]))$/,
            date: /^\d{4}-(?:0[0-9]{1}|1[0-2]{1})-(3[01]|0[1-9]|[12][0-9])$/,
            time: /^(2[0-4]|[01][0-9]):([0-5][0-9]):(60|[0-5][0-9])$/,
            email: /^(?:[\w\!\#\$\%\&\'\*\+\-\/\=\?\^\`\{\|\}\~]+\.)*[\w\!\#\$\%\&\'\*\+\-\/\=\?\^\`\{\|\}\~]+@(?:(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9\-](?!\.)){0,61}[a-zA-Z0-9]?\.)+[a-zA-Z0-9](?:[a-zA-Z0-9\-](?!$)){0,61}[a-zA-Z0-9]?)|(?:\[(?:(?:[01]?\d{1,2}|2[0-4]\d|25[0-5])\.){3}(?:[01]?\d{1,2}|2[0-4]\d|25[0-5])\]))$/,
            "ip-address": /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
            ipv6: /^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$/,
            uri: /^[a-zA-Z][a-zA-Z0-9+-.]*:[^\s]*$/,
            "uri-reference": /^(((([A-Za-z][+\-.0-9A-Za-z]*(:%[0-9A-Fa-f]{2}|:[!$&-.0-;=?-Z_a-z~]|[/?])|\?)(%[0-9A-Fa-f]{2}|[!$&-;=?-Z_a-z~])*|([A-Za-z][+\-.0-9A-Za-z]*:?)?)|([A-Za-z][+\-.0-9A-Za-z]*:)?\/((%[0-9A-Fa-f]{2}|\/((%[0-9A-Fa-f]{2}|[!$&-.0-9;=A-Z_a-z~])+|(\[(([Vv][0-9A-Fa-f]+\.[!$&-.0-;=A-Z_a-z~]+)?|[.0-:A-Fa-f]+)\])?)(:\d*)?[/?]|[!$&-.0-;=?-Z_a-z~])(%[0-9A-Fa-f]{2}|[!$&-;=?-Z_a-z~])*|(\/((%[0-9A-Fa-f]{2}|[!$&-.0-9;=A-Z_a-z~])+|(\[(([Vv][0-9A-Fa-f]+\.[!$&-.0-;=A-Z_a-z~]+)?|[.0-:A-Fa-f]+)\])?)(:\d*)?)?))#(%[0-9A-Fa-f]{2}|[!$&-;=?-Z_a-z~])*|(([A-Za-z][+\-.0-9A-Za-z]*)?%[0-9A-Fa-f]{2}|[!$&-.0-9;=@_~]|[A-Za-z][+\-.0-9A-Za-z]*[!$&-*,;=@_~])(%[0-9A-Fa-f]{2}|[!$&-.0-9;=@-Z_a-z~])*((([/?](%[0-9A-Fa-f]{2}|[!$&-;=?-Z_a-z~])*)?#|[/?])(%[0-9A-Fa-f]{2}|[!$&-;=?-Z_a-z~])*)?|([A-Za-z][+\-.0-9A-Za-z]*(:%[0-9A-Fa-f]{2}|:[!$&-.0-;=?-Z_a-z~]|[/?])|\?)(%[0-9A-Fa-f]{2}|[!$&-;=?-Z_a-z~])*|([A-Za-z][+\-.0-9A-Za-z]*:)?\/((%[0-9A-Fa-f]{2}|\/((%[0-9A-Fa-f]{2}|[!$&-.0-9;=A-Z_a-z~])+|(\[(([Vv][0-9A-Fa-f]+\.[!$&-.0-;=A-Z_a-z~]+)?|[.0-:A-Fa-f]+)\])?)(:\d*)?[/?]|[!$&-.0-;=?-Z_a-z~])(%[0-9A-Fa-f]{2}|[!$&-;=?-Z_a-z~])*|\/((%[0-9A-Fa-f]{2}|[!$&-.0-9;=A-Z_a-z~])+(:\d*)?|(\[(([Vv][0-9A-Fa-f]+\.[!$&-.0-;=A-Z_a-z~]+)?|[.0-:A-Fa-f]+)\])?:\d*|\[(([Vv][0-9A-Fa-f]+\.[!$&-.0-;=A-Z_a-z~]+)?|[.0-:A-Fa-f]+)\])?)?|[A-Za-z][+\-.0-9A-Za-z]*:?)?$/,
            color: /^(#?([0-9A-Fa-f]{3}){1,2}\b|aqua|black|blue|fuchsia|gray|green|lime|maroon|navy|olive|orange|purple|red|silver|teal|white|yellow|(rgb\(\s*\b([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\b\s*,\s*\b([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\b\s*,\s*\b([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\b\s*\))|(rgb\(\s*(\d?\d%|100%)+\s*,\s*(\d?\d%|100%)+\s*,\s*(\d?\d%|100%)+\s*\)))$/,
            hostname: /^(?=.{1,255}$)[0-9A-Za-z](?:(?:[0-9A-Za-z]|-){0,61}[0-9A-Za-z])?(?:\.[0-9A-Za-z](?:(?:[0-9A-Za-z]|-){0,61}[0-9A-Za-z])?)*\.?$/,
            "host-name": /^(?=.{1,255}$)[0-9A-Za-z](?:(?:[0-9A-Za-z]|-){0,61}[0-9A-Za-z])?(?:\.[0-9A-Za-z](?:(?:[0-9A-Za-z]|-){0,61}[0-9A-Za-z])?)*\.?$/,
            alpha: /^[a-zA-Z]+$/,
            alphanumeric: /^[a-zA-Z0-9]+$/,
            "utc-millisec": function(input) {
                return typeof input === "string" && parseFloat(input) === parseInt(input, 10) && !isNaN(input);
            },
            regex: function(input) {
                var result = true;
                try {
                    new RegExp(input);
                } catch (e) {
                    result = false;
                }
                return result;
            },
            style: /\s*(.+?):\s*([^;]+);?/,
            phone: /^\+(?:[0-9] ?){6,14}[0-9]$/
        };
        FORMAT_REGEXPS.regexp = FORMAT_REGEXPS.regex;
        FORMAT_REGEXPS.pattern = FORMAT_REGEXPS.regex;
        FORMAT_REGEXPS.ipv4 = FORMAT_REGEXPS["ip-address"];
        exports.isFormat = function isFormat(input, format, validator) {
            if (typeof input === "string" && FORMAT_REGEXPS[format] !== undefined) {
                if (FORMAT_REGEXPS[format] instanceof RegExp) {
                    return FORMAT_REGEXPS[format].test(input);
                }
                if (typeof FORMAT_REGEXPS[format] === "function") {
                    return FORMAT_REGEXPS[format](input);
                }
            } else if (validator && validator.customFormats && typeof validator.customFormats[format] === "function") {
                return validator.customFormats[format](input);
            }
            return true;
        };
        var makeSuffix = exports.makeSuffix = function makeSuffix(key) {
            key = key.toString();
            if (!key.match(/[.\s\[\]]/) && !key.match(/^[\d]/)) {
                return "." + key;
            }
            if (key.match(/^\d+$/)) {
                return "[" + key + "]";
            }
            return "[" + JSON.stringify(key) + "]";
        };
        exports.deepCompareStrict = function deepCompareStrict(a, b) {
            if (typeof a !== typeof b) {
                return false;
            }
            if (Array.isArray(a)) {
                if (!Array.isArray(b)) {
                    return false;
                }
                if (a.length !== b.length) {
                    return false;
                }
                return a.every((function(v, i) {
                    return deepCompareStrict(a[i], b[i]);
                }));
            }
            if (typeof a === "object") {
                if (!a || !b) {
                    return a === b;
                }
                var aKeys = Object.keys(a);
                var bKeys = Object.keys(b);
                if (aKeys.length !== bKeys.length) {
                    return false;
                }
                return aKeys.every((function(v) {
                    return deepCompareStrict(a[v], b[v]);
                }));
            }
            return a === b;
        };
        function deepMerger(target, dst, e, i) {
            if (typeof e === "object") {
                dst[i] = deepMerge(target[i], e);
            } else {
                if (target.indexOf(e) === -1) {
                    dst.push(e);
                }
            }
        }
        function copyist(src, dst, key) {
            dst[key] = src[key];
        }
        function copyistWithDeepMerge(target, src, dst, key) {
            if (typeof src[key] !== "object" || !src[key]) {
                dst[key] = src[key];
            } else {
                if (!target[key]) {
                    dst[key] = src[key];
                } else {
                    dst[key] = deepMerge(target[key], src[key]);
                }
            }
        }
        function deepMerge(target, src) {
            var array = Array.isArray(src);
            var dst = array && [] || {};
            if (array) {
                target = target || [];
                dst = dst.concat(target);
                src.forEach(deepMerger.bind(null, target, dst));
            } else {
                if (target && typeof target === "object") {
                    Object.keys(target).forEach(copyist.bind(null, target, dst));
                }
                Object.keys(src).forEach(copyistWithDeepMerge.bind(null, target, src, dst));
            }
            return dst;
        }
        module.exports.deepMerge = deepMerge;
        exports.objectGetPath = function objectGetPath(o, s) {
            var parts = s.split("/").slice(1);
            var k;
            while (typeof (k = parts.shift()) == "string") {
                var n = decodeURIComponent(k.replace(/~0/, "~").replace(/~1/g, "/"));
                if (!(n in o)) return;
                o = o[n];
            }
            return o;
        };
        function pathEncoder(v) {
            return "/" + encodeURIComponent(v).replace(/~/g, "%7E");
        }
        exports.encodePath = function encodePointer(a) {
            return a.map(pathEncoder).join("");
        };
        exports.getDecimalPlaces = function getDecimalPlaces(number) {
            var decimalPlaces = 0;
            if (isNaN(number)) return decimalPlaces;
            if (typeof number !== "number") {
                number = Number(number);
            }
            var parts = number.toString().split("e");
            if (parts.length === 2) {
                if (parts[1][0] !== "-") {
                    return decimalPlaces;
                } else {
                    decimalPlaces = Number(parts[1].slice(1));
                }
            }
            var decimalParts = parts[0].split(".");
            if (decimalParts.length === 2) {
                decimalPlaces += decimalParts[1].length;
            }
            return decimalPlaces;
        };
        exports.isSchema = function isSchema(val) {
            return typeof val === "object" && val || typeof val === "boolean";
        };
    },
    6872: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        var Validator = module.exports.Validator = __webpack_require__(8175);
        module.exports.ValidatorResult = __webpack_require__(6495).ValidatorResult;
        module.exports.ValidatorResultError = __webpack_require__(6495).ValidatorResultError;
        module.exports.ValidationError = __webpack_require__(6495).ValidationError;
        module.exports.SchemaError = __webpack_require__(6495).SchemaError;
        module.exports.SchemaScanResult = __webpack_require__(4496).SchemaScanResult;
        module.exports.scan = __webpack_require__(4496).scan;
        module.exports.validate = function(instance, schema, options) {
            var v = new Validator;
            return v.validate(instance, schema, options);
        };
    },
    4496: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        var urilib = __webpack_require__(8835);
        var helpers = __webpack_require__(6495);
        module.exports.SchemaScanResult = SchemaScanResult;
        function SchemaScanResult(found, ref) {
            this.id = found;
            this.ref = ref;
        }
        module.exports.scan = function scan(base, schema) {
            function scanSchema(baseuri, schema) {
                if (!schema || typeof schema != "object") return;
                if (schema.$ref) {
                    var resolvedUri = urilib.resolve(baseuri, schema.$ref);
                    ref[resolvedUri] = ref[resolvedUri] ? ref[resolvedUri] + 1 : 0;
                    return;
                }
                var id = schema.$id || schema.id;
                var ourBase = id ? urilib.resolve(baseuri, id) : baseuri;
                if (ourBase) {
                    if (ourBase.indexOf("#") < 0) ourBase += "#";
                    if (found[ourBase]) {
                        if (!helpers.deepCompareStrict(found[ourBase], schema)) {
                            throw new Error("Schema <" + ourBase + "> already exists with different definition");
                        }
                        return found[ourBase];
                    }
                    found[ourBase] = schema;
                    if (ourBase[ourBase.length - 1] == "#") {
                        found[ourBase.substring(0, ourBase.length - 1)] = schema;
                    }
                }
                scanArray(ourBase + "/items", Array.isArray(schema.items) ? schema.items : [ schema.items ]);
                scanArray(ourBase + "/extends", Array.isArray(schema.extends) ? schema.extends : [ schema.extends ]);
                scanSchema(ourBase + "/additionalItems", schema.additionalItems);
                scanObject(ourBase + "/properties", schema.properties);
                scanSchema(ourBase + "/additionalProperties", schema.additionalProperties);
                scanObject(ourBase + "/definitions", schema.definitions);
                scanObject(ourBase + "/patternProperties", schema.patternProperties);
                scanObject(ourBase + "/dependencies", schema.dependencies);
                scanArray(ourBase + "/disallow", schema.disallow);
                scanArray(ourBase + "/allOf", schema.allOf);
                scanArray(ourBase + "/anyOf", schema.anyOf);
                scanArray(ourBase + "/oneOf", schema.oneOf);
                scanSchema(ourBase + "/not", schema.not);
            }
            function scanArray(baseuri, schemas) {
                if (!Array.isArray(schemas)) return;
                for (var i = 0; i < schemas.length; i++) {
                    scanSchema(baseuri + "/" + i, schemas[i]);
                }
            }
            function scanObject(baseuri, schemas) {
                if (!schemas || typeof schemas != "object") return;
                for (var p in schemas) {
                    scanSchema(baseuri + "/" + p, schemas[p]);
                }
            }
            var found = {};
            var ref = {};
            scanSchema(base, schema);
            return new SchemaScanResult(found, ref);
        };
    },
    8175: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        var urilib = __webpack_require__(8835);
        var attribute = __webpack_require__(6120);
        var helpers = __webpack_require__(6495);
        var scanSchema = __webpack_require__(4496).scan;
        var ValidatorResult = helpers.ValidatorResult;
        var ValidatorResultError = helpers.ValidatorResultError;
        var SchemaError = helpers.SchemaError;
        var SchemaContext = helpers.SchemaContext;
        var anonymousBase = "/";
        var Validator = function Validator() {
            this.customFormats = Object.create(Validator.prototype.customFormats);
            this.schemas = {};
            this.unresolvedRefs = [];
            this.types = Object.create(types);
            this.attributes = Object.create(attribute.validators);
        };
        Validator.prototype.customFormats = {};
        Validator.prototype.schemas = null;
        Validator.prototype.types = null;
        Validator.prototype.attributes = null;
        Validator.prototype.unresolvedRefs = null;
        Validator.prototype.addSchema = function addSchema(schema, base) {
            var self = this;
            if (!schema) {
                return null;
            }
            var scan = scanSchema(base || anonymousBase, schema);
            var ourUri = base || schema.$id || schema.id;
            for (var uri in scan.id) {
                this.schemas[uri] = scan.id[uri];
            }
            for (var uri in scan.ref) {
                this.unresolvedRefs.push(uri);
            }
            this.unresolvedRefs = this.unresolvedRefs.filter((function(uri) {
                return typeof self.schemas[uri] === "undefined";
            }));
            return this.schemas[ourUri];
        };
        Validator.prototype.addSubSchemaArray = function addSubSchemaArray(baseuri, schemas) {
            if (!Array.isArray(schemas)) return;
            for (var i = 0; i < schemas.length; i++) {
                this.addSubSchema(baseuri, schemas[i]);
            }
        };
        Validator.prototype.addSubSchemaObject = function addSubSchemaArray(baseuri, schemas) {
            if (!schemas || typeof schemas != "object") return;
            for (var p in schemas) {
                this.addSubSchema(baseuri, schemas[p]);
            }
        };
        Validator.prototype.setSchemas = function setSchemas(schemas) {
            this.schemas = schemas;
        };
        Validator.prototype.getSchema = function getSchema(urn) {
            return this.schemas[urn];
        };
        Validator.prototype.validate = function validate(instance, schema, options, ctx) {
            if (typeof schema !== "boolean" && typeof schema !== "object" || schema === null) {
                throw new SchemaError("Expected `schema` to be an object or boolean");
            }
            if (!options) {
                options = {};
            }
            var id = schema.$id || schema.id;
            var base = urilib.resolve(options.base || anonymousBase, id || "");
            if (!ctx) {
                ctx = new SchemaContext(schema, options, [], base, Object.create(this.schemas));
                if (!ctx.schemas[base]) {
                    ctx.schemas[base] = schema;
                }
                var found = scanSchema(base, schema);
                for (var n in found.id) {
                    var sch = found.id[n];
                    ctx.schemas[n] = sch;
                }
            }
            if (options.required && instance === undefined) {
                var result = new ValidatorResult(instance, schema, options, ctx);
                result.addError("is required, but is undefined");
                return result;
            }
            var result = this.validateSchema(instance, schema, options, ctx);
            if (!result) {
                throw new Error("Result undefined");
            } else if (options.throwAll && result.errors.length) {
                throw new ValidatorResultError(result);
            }
            return result;
        };
        function shouldResolve(schema) {
            var ref = typeof schema === "string" ? schema : schema.$ref;
            if (typeof ref == "string") return ref;
            return false;
        }
        Validator.prototype.validateSchema = function validateSchema(instance, schema, options, ctx) {
            var result = new ValidatorResult(instance, schema, options, ctx);
            if (typeof schema === "boolean") {
                if (schema === true) {
                    schema = {};
                } else if (schema === false) {
                    schema = {
                        type: []
                    };
                }
            } else if (!schema) {
                throw new Error("schema is undefined");
            }
            if (schema["extends"]) {
                if (Array.isArray(schema["extends"])) {
                    var schemaobj = {
                        schema,
                        ctx
                    };
                    schema["extends"].forEach(this.schemaTraverser.bind(this, schemaobj));
                    schema = schemaobj.schema;
                    schemaobj.schema = null;
                    schemaobj.ctx = null;
                    schemaobj = null;
                } else {
                    schema = helpers.deepMerge(schema, this.superResolve(schema["extends"], ctx));
                }
            }
            var switchSchema = shouldResolve(schema);
            if (switchSchema) {
                var resolved = this.resolve(schema, switchSchema, ctx);
                var subctx = new SchemaContext(resolved.subschema, options, ctx.path, resolved.switchSchema, ctx.schemas);
                return this.validateSchema(instance, resolved.subschema, options, subctx);
            }
            var skipAttributes = options && options.skipAttributes || [];
            for (var key in schema) {
                if (!attribute.ignoreProperties[key] && skipAttributes.indexOf(key) < 0) {
                    var validatorErr = null;
                    var validator = this.attributes[key];
                    if (validator) {
                        validatorErr = validator.call(this, instance, schema, options, ctx);
                    } else if (options.allowUnknownAttributes === false) {
                        throw new SchemaError("Unsupported attribute: " + key, schema);
                    }
                    if (validatorErr) {
                        result.importErrors(validatorErr);
                    }
                }
            }
            if (typeof options.rewrite == "function") {
                var value = options.rewrite.call(this, instance, schema, options, ctx);
                result.instance = value;
            }
            return result;
        };
        Validator.prototype.schemaTraverser = function schemaTraverser(schemaobj, s) {
            schemaobj.schema = helpers.deepMerge(schemaobj.schema, this.superResolve(s, schemaobj.ctx));
        };
        Validator.prototype.superResolve = function superResolve(schema, ctx) {
            var ref = shouldResolve(schema);
            if (ref) {
                return this.resolve(schema, ref, ctx).subschema;
            }
            return schema;
        };
        Validator.prototype.resolve = function resolve(schema, switchSchema, ctx) {
            switchSchema = ctx.resolve(switchSchema);
            if (ctx.schemas[switchSchema]) {
                return {
                    subschema: ctx.schemas[switchSchema],
                    switchSchema
                };
            }
            var parsed = urilib.parse(switchSchema);
            var fragment = parsed && parsed.hash;
            var document = fragment && fragment.length && switchSchema.substr(0, switchSchema.length - fragment.length);
            if (!document || !ctx.schemas[document]) {
                throw new SchemaError("no such schema <" + switchSchema + ">", schema);
            }
            var subschema = helpers.objectGetPath(ctx.schemas[document], fragment.substr(1));
            if (subschema === undefined) {
                throw new SchemaError("no such schema " + fragment + " located in <" + document + ">", schema);
            }
            return {
                subschema,
                switchSchema
            };
        };
        Validator.prototype.testType = function validateType(instance, schema, options, ctx, type) {
            if (type === undefined) {
                return;
            } else if (type === null) {
                throw new SchemaError('Unexpected null in "type" keyword');
            }
            if (typeof this.types[type] == "function") {
                return this.types[type].call(this, instance);
            }
            if (type && typeof type == "object") {
                var res = this.validateSchema(instance, type, options, ctx);
                return res === undefined || !(res && res.errors.length);
            }
            return true;
        };
        var types = Validator.prototype.types = {};
        types.string = function testString(instance) {
            return typeof instance == "string";
        };
        types.number = function testNumber(instance) {
            return typeof instance == "number" && isFinite(instance);
        };
        types.integer = function testInteger(instance) {
            return typeof instance == "number" && instance % 1 === 0;
        };
        types.boolean = function testBoolean(instance) {
            return typeof instance == "boolean";
        };
        types.array = function testArray(instance) {
            return Array.isArray(instance);
        };
        types["null"] = function testNull(instance) {
            return instance === null;
        };
        types.date = function testDate(instance) {
            return instance instanceof Date;
        };
        types.any = function testAny(instance) {
            return true;
        };
        types.object = function testObject(instance) {
            return instance && typeof instance === "object" && !Array.isArray(instance) && !(instance instanceof Date);
        };
        module.exports = Validator;
    },
    2253: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const EE = __webpack_require__(8614);
        const Stream = __webpack_require__(2413);
        const Yallist = __webpack_require__(1455);
        const SD = __webpack_require__(4304).StringDecoder;
        const EOF = Symbol("EOF");
        const MAYBE_EMIT_END = Symbol("maybeEmitEnd");
        const EMITTED_END = Symbol("emittedEnd");
        const EMITTING_END = Symbol("emittingEnd");
        const CLOSED = Symbol("closed");
        const READ = Symbol("read");
        const FLUSH = Symbol("flush");
        const FLUSHCHUNK = Symbol("flushChunk");
        const ENCODING = Symbol("encoding");
        const DECODER = Symbol("decoder");
        const FLOWING = Symbol("flowing");
        const PAUSED = Symbol("paused");
        const RESUME = Symbol("resume");
        const BUFFERLENGTH = Symbol("bufferLength");
        const BUFFERPUSH = Symbol("bufferPush");
        const BUFFERSHIFT = Symbol("bufferShift");
        const OBJECTMODE = Symbol("objectMode");
        const DESTROYED = Symbol("destroyed");
        const doIter = global._MP_NO_ITERATOR_SYMBOLS_ !== "1";
        const ASYNCITERATOR = doIter && Symbol.asyncIterator || Symbol("asyncIterator not implemented");
        const ITERATOR = doIter && Symbol.iterator || Symbol("iterator not implemented");
        const isEndish = ev => ev === "end" || ev === "finish" || ev === "prefinish";
        const isArrayBuffer = b => b instanceof ArrayBuffer || typeof b === "object" && b.constructor && b.constructor.name === "ArrayBuffer" && b.byteLength >= 0;
        const isArrayBufferView = b => !Buffer.isBuffer(b) && ArrayBuffer.isView(b);
        module.exports = class Minipass extends Stream {
            constructor(options) {
                super();
                this[FLOWING] = false;
                this[PAUSED] = false;
                this.pipes = new Yallist;
                this.buffer = new Yallist;
                this[OBJECTMODE] = options && options.objectMode || false;
                if (this[OBJECTMODE]) this[ENCODING] = null; else this[ENCODING] = options && options.encoding || null;
                if (this[ENCODING] === "buffer") this[ENCODING] = null;
                this[DECODER] = this[ENCODING] ? new SD(this[ENCODING]) : null;
                this[EOF] = false;
                this[EMITTED_END] = false;
                this[EMITTING_END] = false;
                this[CLOSED] = false;
                this.writable = true;
                this.readable = true;
                this[BUFFERLENGTH] = 0;
                this[DESTROYED] = false;
            }
            get bufferLength() {
                return this[BUFFERLENGTH];
            }
            get encoding() {
                return this[ENCODING];
            }
            set encoding(enc) {
                if (this[OBJECTMODE]) throw new Error("cannot set encoding in objectMode");
                if (this[ENCODING] && enc !== this[ENCODING] && (this[DECODER] && this[DECODER].lastNeed || this[BUFFERLENGTH])) throw new Error("cannot change encoding");
                if (this[ENCODING] !== enc) {
                    this[DECODER] = enc ? new SD(enc) : null;
                    if (this.buffer.length) this.buffer = this.buffer.map((chunk => this[DECODER].write(chunk)));
                }
                this[ENCODING] = enc;
            }
            setEncoding(enc) {
                this.encoding = enc;
            }
            get objectMode() {
                return this[OBJECTMODE];
            }
            set objectMode(om) {
                this[OBJECTMODE] = this[OBJECTMODE] || !!om;
            }
            write(chunk, encoding, cb) {
                if (this[EOF]) throw new Error("write after end");
                if (this[DESTROYED]) {
                    this.emit("error", Object.assign(new Error("Cannot call write after a stream was destroyed"), {
                        code: "ERR_STREAM_DESTROYED"
                    }));
                    return true;
                }
                if (typeof encoding === "function") cb = encoding, encoding = "utf8";
                if (!encoding) encoding = "utf8";
                if (!this[OBJECTMODE] && !Buffer.isBuffer(chunk)) {
                    if (isArrayBufferView(chunk)) chunk = Buffer.from(chunk.buffer, chunk.byteOffset, chunk.byteLength); else if (isArrayBuffer(chunk)) chunk = Buffer.from(chunk); else if (typeof chunk !== "string") this.objectMode = true;
                }
                if (!this.objectMode && !chunk.length) {
                    if (this[BUFFERLENGTH] !== 0) this.emit("readable");
                    if (cb) cb();
                    return this.flowing;
                }
                if (typeof chunk === "string" && !this[OBJECTMODE] && !(encoding === this[ENCODING] && !this[DECODER].lastNeed)) {
                    chunk = Buffer.from(chunk, encoding);
                }
                if (Buffer.isBuffer(chunk) && this[ENCODING]) chunk = this[DECODER].write(chunk);
                if (this.flowing) {
                    if (this[BUFFERLENGTH] !== 0) this[FLUSH](true);
                    this.emit("data", chunk);
                } else this[BUFFERPUSH](chunk);
                if (this[BUFFERLENGTH] !== 0) this.emit("readable");
                if (cb) cb();
                return this.flowing;
            }
            read(n) {
                if (this[DESTROYED]) return null;
                try {
                    if (this[BUFFERLENGTH] === 0 || n === 0 || n > this[BUFFERLENGTH]) return null;
                    if (this[OBJECTMODE]) n = null;
                    if (this.buffer.length > 1 && !this[OBJECTMODE]) {
                        if (this.encoding) this.buffer = new Yallist([ Array.from(this.buffer).join("") ]); else this.buffer = new Yallist([ Buffer.concat(Array.from(this.buffer), this[BUFFERLENGTH]) ]);
                    }
                    return this[READ](n || null, this.buffer.head.value);
                } finally {
                    this[MAYBE_EMIT_END]();
                }
            }
            [READ](n, chunk) {
                if (n === chunk.length || n === null) this[BUFFERSHIFT](); else {
                    this.buffer.head.value = chunk.slice(n);
                    chunk = chunk.slice(0, n);
                    this[BUFFERLENGTH] -= n;
                }
                this.emit("data", chunk);
                if (!this.buffer.length && !this[EOF]) this.emit("drain");
                return chunk;
            }
            end(chunk, encoding, cb) {
                if (typeof chunk === "function") cb = chunk, chunk = null;
                if (typeof encoding === "function") cb = encoding, encoding = "utf8";
                if (chunk) this.write(chunk, encoding);
                if (cb) this.once("end", cb);
                this[EOF] = true;
                this.writable = false;
                if (this.flowing || !this[PAUSED]) this[MAYBE_EMIT_END]();
                return this;
            }
            [RESUME]() {
                if (this[DESTROYED]) return;
                this[PAUSED] = false;
                this[FLOWING] = true;
                this.emit("resume");
                if (this.buffer.length) this[FLUSH](); else if (this[EOF]) this[MAYBE_EMIT_END](); else this.emit("drain");
            }
            resume() {
                return this[RESUME]();
            }
            pause() {
                this[FLOWING] = false;
                this[PAUSED] = true;
            }
            get destroyed() {
                return this[DESTROYED];
            }
            get flowing() {
                return this[FLOWING];
            }
            get paused() {
                return this[PAUSED];
            }
            [BUFFERPUSH](chunk) {
                if (this[OBJECTMODE]) this[BUFFERLENGTH] += 1; else this[BUFFERLENGTH] += chunk.length;
                return this.buffer.push(chunk);
            }
            [BUFFERSHIFT]() {
                if (this.buffer.length) {
                    if (this[OBJECTMODE]) this[BUFFERLENGTH] -= 1; else this[BUFFERLENGTH] -= this.buffer.head.value.length;
                }
                return this.buffer.shift();
            }
            [FLUSH](noDrain) {
                do {} while (this[FLUSHCHUNK](this[BUFFERSHIFT]()));
                if (!noDrain && !this.buffer.length && !this[EOF]) this.emit("drain");
            }
            [FLUSHCHUNK](chunk) {
                return chunk ? (this.emit("data", chunk), this.flowing) : false;
            }
            pipe(dest, opts) {
                if (this[DESTROYED]) return;
                const ended = this[EMITTED_END];
                opts = opts || {};
                if (dest === process.stdout || dest === process.stderr) opts.end = false; else opts.end = opts.end !== false;
                const p = {
                    dest,
                    opts,
                    ondrain: _ => this[RESUME]()
                };
                this.pipes.push(p);
                dest.on("drain", p.ondrain);
                this[RESUME]();
                if (ended && p.opts.end) p.dest.end();
                return dest;
            }
            addListener(ev, fn) {
                return this.on(ev, fn);
            }
            on(ev, fn) {
                try {
                    return super.on(ev, fn);
                } finally {
                    if (ev === "data" && !this.pipes.length && !this.flowing) this[RESUME](); else if (isEndish(ev) && this[EMITTED_END]) {
                        super.emit(ev);
                        this.removeAllListeners(ev);
                    }
                }
            }
            get emittedEnd() {
                return this[EMITTED_END];
            }
            [MAYBE_EMIT_END]() {
                if (!this[EMITTING_END] && !this[EMITTED_END] && !this[DESTROYED] && this.buffer.length === 0 && this[EOF]) {
                    this[EMITTING_END] = true;
                    this.emit("end");
                    this.emit("prefinish");
                    this.emit("finish");
                    if (this[CLOSED]) this.emit("close");
                    this[EMITTING_END] = false;
                }
            }
            emit(ev, data) {
                if (ev !== "error" && ev !== "close" && ev !== DESTROYED && this[DESTROYED]) return; else if (ev === "data") {
                    if (!data) return;
                    if (this.pipes.length) this.pipes.forEach((p => p.dest.write(data) === false && this.pause()));
                } else if (ev === "end") {
                    if (this[EMITTED_END] === true) return;
                    this[EMITTED_END] = true;
                    this.readable = false;
                    if (this[DECODER]) {
                        data = this[DECODER].end();
                        if (data) {
                            this.pipes.forEach((p => p.dest.write(data)));
                            super.emit("data", data);
                        }
                    }
                    this.pipes.forEach((p => {
                        p.dest.removeListener("drain", p.ondrain);
                        if (p.opts.end) p.dest.end();
                    }));
                } else if (ev === "close") {
                    this[CLOSED] = true;
                    if (!this[EMITTED_END] && !this[DESTROYED]) return;
                }
                const args = new Array(arguments.length);
                args[0] = ev;
                args[1] = data;
                if (arguments.length > 2) {
                    for (let i = 2; i < arguments.length; i++) {
                        args[i] = arguments[i];
                    }
                }
                try {
                    return super.emit.apply(this, args);
                } finally {
                    if (!isEndish(ev)) this[MAYBE_EMIT_END](); else this.removeAllListeners(ev);
                }
            }
            collect() {
                const buf = [];
                if (!this[OBJECTMODE]) buf.dataLength = 0;
                const p = this.promise();
                this.on("data", (c => {
                    buf.push(c);
                    if (!this[OBJECTMODE]) buf.dataLength += c.length;
                }));
                return p.then((() => buf));
            }
            concat() {
                return this[OBJECTMODE] ? Promise.reject(new Error("cannot concat in objectMode")) : this.collect().then((buf => this[OBJECTMODE] ? Promise.reject(new Error("cannot concat in objectMode")) : this[ENCODING] ? buf.join("") : Buffer.concat(buf, buf.dataLength)));
            }
            promise() {
                return new Promise(((resolve, reject) => {
                    this.on(DESTROYED, (() => reject(new Error("stream destroyed"))));
                    this.on("end", (() => resolve()));
                    this.on("error", (er => reject(er)));
                }));
            }
            [ASYNCITERATOR]() {
                const next = () => {
                    const res = this.read();
                    if (res !== null) return Promise.resolve({
                        done: false,
                        value: res
                    });
                    if (this[EOF]) return Promise.resolve({
                        done: true
                    });
                    let resolve = null;
                    let reject = null;
                    const onerr = er => {
                        this.removeListener("data", ondata);
                        this.removeListener("end", onend);
                        reject(er);
                    };
                    const ondata = value => {
                        this.removeListener("error", onerr);
                        this.removeListener("end", onend);
                        this.pause();
                        resolve({
                            value,
                            done: !!this[EOF]
                        });
                    };
                    const onend = () => {
                        this.removeListener("error", onerr);
                        this.removeListener("data", ondata);
                        resolve({
                            done: true
                        });
                    };
                    const ondestroy = () => onerr(new Error("stream destroyed"));
                    return new Promise(((res, rej) => {
                        reject = rej;
                        resolve = res;
                        this.once(DESTROYED, ondestroy);
                        this.once("error", onerr);
                        this.once("end", onend);
                        this.once("data", ondata);
                    }));
                };
                return {
                    next
                };
            }
            [ITERATOR]() {
                const next = () => {
                    const value = this.read();
                    const done = value === null;
                    return {
                        value,
                        done
                    };
                };
                return {
                    next
                };
            }
            destroy(er) {
                if (this[DESTROYED]) {
                    if (er) this.emit("error", er); else this.emit(DESTROYED);
                    return this;
                }
                this[DESTROYED] = true;
                this.buffer = new Yallist;
                this[BUFFERLENGTH] = 0;
                if (typeof this.close === "function" && !this[CLOSED]) this.close();
                if (er) this.emit("error", er); else this.emit(DESTROYED);
                return this;
            }
            static isStream(s) {
                return !!s && (s instanceof Minipass || s instanceof Stream || s instanceof EE && (typeof s.pipe === "function" || typeof s.write === "function" && typeof s.end === "function"));
            }
        };
    },
    8597: (module, __unused_webpack_exports, __webpack_require__) => {
        const realZlibConstants = __webpack_require__(8761).constants || {
            ZLIB_VERNUM: 4736
        };
        module.exports = Object.freeze(Object.assign(Object.create(null), {
            Z_NO_FLUSH: 0,
            Z_PARTIAL_FLUSH: 1,
            Z_SYNC_FLUSH: 2,
            Z_FULL_FLUSH: 3,
            Z_FINISH: 4,
            Z_BLOCK: 5,
            Z_OK: 0,
            Z_STREAM_END: 1,
            Z_NEED_DICT: 2,
            Z_ERRNO: -1,
            Z_STREAM_ERROR: -2,
            Z_DATA_ERROR: -3,
            Z_MEM_ERROR: -4,
            Z_BUF_ERROR: -5,
            Z_VERSION_ERROR: -6,
            Z_NO_COMPRESSION: 0,
            Z_BEST_SPEED: 1,
            Z_BEST_COMPRESSION: 9,
            Z_DEFAULT_COMPRESSION: -1,
            Z_FILTERED: 1,
            Z_HUFFMAN_ONLY: 2,
            Z_RLE: 3,
            Z_FIXED: 4,
            Z_DEFAULT_STRATEGY: 0,
            DEFLATE: 1,
            INFLATE: 2,
            GZIP: 3,
            GUNZIP: 4,
            DEFLATERAW: 5,
            INFLATERAW: 6,
            UNZIP: 7,
            BROTLI_DECODE: 8,
            BROTLI_ENCODE: 9,
            Z_MIN_WINDOWBITS: 8,
            Z_MAX_WINDOWBITS: 15,
            Z_DEFAULT_WINDOWBITS: 15,
            Z_MIN_CHUNK: 64,
            Z_MAX_CHUNK: Infinity,
            Z_DEFAULT_CHUNK: 16384,
            Z_MIN_MEMLEVEL: 1,
            Z_MAX_MEMLEVEL: 9,
            Z_DEFAULT_MEMLEVEL: 8,
            Z_MIN_LEVEL: -1,
            Z_MAX_LEVEL: 9,
            Z_DEFAULT_LEVEL: -1,
            BROTLI_OPERATION_PROCESS: 0,
            BROTLI_OPERATION_FLUSH: 1,
            BROTLI_OPERATION_FINISH: 2,
            BROTLI_OPERATION_EMIT_METADATA: 3,
            BROTLI_MODE_GENERIC: 0,
            BROTLI_MODE_TEXT: 1,
            BROTLI_MODE_FONT: 2,
            BROTLI_DEFAULT_MODE: 0,
            BROTLI_MIN_QUALITY: 0,
            BROTLI_MAX_QUALITY: 11,
            BROTLI_DEFAULT_QUALITY: 11,
            BROTLI_MIN_WINDOW_BITS: 10,
            BROTLI_MAX_WINDOW_BITS: 24,
            BROTLI_LARGE_MAX_WINDOW_BITS: 30,
            BROTLI_DEFAULT_WINDOW: 22,
            BROTLI_MIN_INPUT_BLOCK_BITS: 16,
            BROTLI_MAX_INPUT_BLOCK_BITS: 24,
            BROTLI_PARAM_MODE: 0,
            BROTLI_PARAM_QUALITY: 1,
            BROTLI_PARAM_LGWIN: 2,
            BROTLI_PARAM_LGBLOCK: 3,
            BROTLI_PARAM_DISABLE_LITERAL_CONTEXT_MODELING: 4,
            BROTLI_PARAM_SIZE_HINT: 5,
            BROTLI_PARAM_LARGE_WINDOW: 6,
            BROTLI_PARAM_NPOSTFIX: 7,
            BROTLI_PARAM_NDIRECT: 8,
            BROTLI_DECODER_RESULT_ERROR: 0,
            BROTLI_DECODER_RESULT_SUCCESS: 1,
            BROTLI_DECODER_RESULT_NEEDS_MORE_INPUT: 2,
            BROTLI_DECODER_RESULT_NEEDS_MORE_OUTPUT: 3,
            BROTLI_DECODER_PARAM_DISABLE_RING_BUFFER_REALLOCATION: 0,
            BROTLI_DECODER_PARAM_LARGE_WINDOW: 1,
            BROTLI_DECODER_NO_ERROR: 0,
            BROTLI_DECODER_SUCCESS: 1,
            BROTLI_DECODER_NEEDS_MORE_INPUT: 2,
            BROTLI_DECODER_NEEDS_MORE_OUTPUT: 3,
            BROTLI_DECODER_ERROR_FORMAT_EXUBERANT_NIBBLE: -1,
            BROTLI_DECODER_ERROR_FORMAT_RESERVED: -2,
            BROTLI_DECODER_ERROR_FORMAT_EXUBERANT_META_NIBBLE: -3,
            BROTLI_DECODER_ERROR_FORMAT_SIMPLE_HUFFMAN_ALPHABET: -4,
            BROTLI_DECODER_ERROR_FORMAT_SIMPLE_HUFFMAN_SAME: -5,
            BROTLI_DECODER_ERROR_FORMAT_CL_SPACE: -6,
            BROTLI_DECODER_ERROR_FORMAT_HUFFMAN_SPACE: -7,
            BROTLI_DECODER_ERROR_FORMAT_CONTEXT_MAP_REPEAT: -8,
            BROTLI_DECODER_ERROR_FORMAT_BLOCK_LENGTH_1: -9,
            BROTLI_DECODER_ERROR_FORMAT_BLOCK_LENGTH_2: -10,
            BROTLI_DECODER_ERROR_FORMAT_TRANSFORM: -11,
            BROTLI_DECODER_ERROR_FORMAT_DICTIONARY: -12,
            BROTLI_DECODER_ERROR_FORMAT_WINDOW_BITS: -13,
            BROTLI_DECODER_ERROR_FORMAT_PADDING_1: -14,
            BROTLI_DECODER_ERROR_FORMAT_PADDING_2: -15,
            BROTLI_DECODER_ERROR_FORMAT_DISTANCE: -16,
            BROTLI_DECODER_ERROR_DICTIONARY_NOT_SET: -19,
            BROTLI_DECODER_ERROR_INVALID_ARGUMENTS: -20,
            BROTLI_DECODER_ERROR_ALLOC_CONTEXT_MODES: -21,
            BROTLI_DECODER_ERROR_ALLOC_TREE_GROUPS: -22,
            BROTLI_DECODER_ERROR_ALLOC_CONTEXT_MAP: -25,
            BROTLI_DECODER_ERROR_ALLOC_RING_BUFFER_1: -26,
            BROTLI_DECODER_ERROR_ALLOC_RING_BUFFER_2: -27,
            BROTLI_DECODER_ERROR_ALLOC_BLOCK_TYPE_TREES: -30,
            BROTLI_DECODER_ERROR_UNREACHABLE: -31
        }, realZlibConstants));
    },
    3704: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        const assert = __webpack_require__(2357);
        const Buffer = __webpack_require__(4293).Buffer;
        const realZlib = __webpack_require__(8761);
        const constants = exports.constants = __webpack_require__(8597);
        const Minipass = __webpack_require__(2253);
        const OriginalBufferConcat = Buffer.concat;
        const _superWrite = Symbol("_superWrite");
        class ZlibError extends Error {
            constructor(err) {
                super("zlib: " + err.message);
                this.code = err.code;
                this.errno = err.errno;
                if (!this.code) this.code = "ZLIB_ERROR";
                this.message = "zlib: " + err.message;
                Error.captureStackTrace(this, this.constructor);
            }
            get name() {
                return "ZlibError";
            }
        }
        const _opts = Symbol("opts");
        const _flushFlag = Symbol("flushFlag");
        const _finishFlushFlag = Symbol("finishFlushFlag");
        const _fullFlushFlag = Symbol("fullFlushFlag");
        const _handle = Symbol("handle");
        const _onError = Symbol("onError");
        const _sawError = Symbol("sawError");
        const _level = Symbol("level");
        const _strategy = Symbol("strategy");
        const _ended = Symbol("ended");
        const _defaultFullFlush = Symbol("_defaultFullFlush");
        class ZlibBase extends Minipass {
            constructor(opts, mode) {
                if (!opts || typeof opts !== "object") throw new TypeError("invalid options for ZlibBase constructor");
                super(opts);
                this[_sawError] = false;
                this[_ended] = false;
                this[_opts] = opts;
                this[_flushFlag] = opts.flush;
                this[_finishFlushFlag] = opts.finishFlush;
                try {
                    this[_handle] = new realZlib[mode](opts);
                } catch (er) {
                    throw new ZlibError(er);
                }
                this[_onError] = err => {
                    if (this[_sawError]) return;
                    this[_sawError] = true;
                    this.close();
                    this.emit("error", err);
                };
                this[_handle].on("error", (er => this[_onError](new ZlibError(er))));
                this.once("end", (() => this.close));
            }
            close() {
                if (this[_handle]) {
                    this[_handle].close();
                    this[_handle] = null;
                    this.emit("close");
                }
            }
            reset() {
                if (!this[_sawError]) {
                    assert(this[_handle], "zlib binding closed");
                    return this[_handle].reset();
                }
            }
            flush(flushFlag) {
                if (this.ended) return;
                if (typeof flushFlag !== "number") flushFlag = this[_fullFlushFlag];
                this.write(Object.assign(Buffer.alloc(0), {
                    [_flushFlag]: flushFlag
                }));
            }
            end(chunk, encoding, cb) {
                if (chunk) this.write(chunk, encoding);
                this.flush(this[_finishFlushFlag]);
                this[_ended] = true;
                return super.end(null, null, cb);
            }
            get ended() {
                return this[_ended];
            }
            write(chunk, encoding, cb) {
                if (typeof encoding === "function") cb = encoding, encoding = "utf8";
                if (typeof chunk === "string") chunk = Buffer.from(chunk, encoding);
                if (this[_sawError]) return;
                assert(this[_handle], "zlib binding closed");
                const nativeHandle = this[_handle]._handle;
                const originalNativeClose = nativeHandle.close;
                nativeHandle.close = () => {};
                const originalClose = this[_handle].close;
                this[_handle].close = () => {};
                Buffer.concat = args => args;
                let result;
                try {
                    const flushFlag = typeof chunk[_flushFlag] === "number" ? chunk[_flushFlag] : this[_flushFlag];
                    result = this[_handle]._processChunk(chunk, flushFlag);
                    Buffer.concat = OriginalBufferConcat;
                } catch (err) {
                    Buffer.concat = OriginalBufferConcat;
                    this[_onError](new ZlibError(err));
                } finally {
                    if (this[_handle]) {
                        this[_handle]._handle = nativeHandle;
                        nativeHandle.close = originalNativeClose;
                        this[_handle].close = originalClose;
                        this[_handle].removeAllListeners("error");
                    }
                }
                if (this[_handle]) this[_handle].on("error", (er => this[_onError](new ZlibError(er))));
                let writeReturn;
                if (result) {
                    if (Array.isArray(result) && result.length > 0) {
                        writeReturn = this[_superWrite](Buffer.from(result[0]));
                        for (let i = 1; i < result.length; i++) {
                            writeReturn = this[_superWrite](result[i]);
                        }
                    } else {
                        writeReturn = this[_superWrite](Buffer.from(result));
                    }
                }
                if (cb) cb();
                return writeReturn;
            }
            [_superWrite](data) {
                return super.write(data);
            }
        }
        class Zlib extends ZlibBase {
            constructor(opts, mode) {
                opts = opts || {};
                opts.flush = opts.flush || constants.Z_NO_FLUSH;
                opts.finishFlush = opts.finishFlush || constants.Z_FINISH;
                super(opts, mode);
                this[_fullFlushFlag] = constants.Z_FULL_FLUSH;
                this[_level] = opts.level;
                this[_strategy] = opts.strategy;
            }
            params(level, strategy) {
                if (this[_sawError]) return;
                if (!this[_handle]) throw new Error("cannot switch params when binding is closed");
                if (!this[_handle].params) throw new Error("not supported in this implementation");
                if (this[_level] !== level || this[_strategy] !== strategy) {
                    this.flush(constants.Z_SYNC_FLUSH);
                    assert(this[_handle], "zlib binding closed");
                    const origFlush = this[_handle].flush;
                    this[_handle].flush = (flushFlag, cb) => {
                        this.flush(flushFlag);
                        cb();
                    };
                    try {
                        this[_handle].params(level, strategy);
                    } finally {
                        this[_handle].flush = origFlush;
                    }
                    if (this[_handle]) {
                        this[_level] = level;
                        this[_strategy] = strategy;
                    }
                }
            }
        }
        class Deflate extends Zlib {
            constructor(opts) {
                super(opts, "Deflate");
            }
        }
        class Inflate extends Zlib {
            constructor(opts) {
                super(opts, "Inflate");
            }
        }
        const _portable = Symbol("_portable");
        class Gzip extends Zlib {
            constructor(opts) {
                super(opts, "Gzip");
                this[_portable] = opts && !!opts.portable;
            }
            [_superWrite](data) {
                if (!this[_portable]) return super[_superWrite](data);
                this[_portable] = false;
                data[9] = 255;
                return super[_superWrite](data);
            }
        }
        class Gunzip extends Zlib {
            constructor(opts) {
                super(opts, "Gunzip");
            }
        }
        class DeflateRaw extends Zlib {
            constructor(opts) {
                super(opts, "DeflateRaw");
            }
        }
        class InflateRaw extends Zlib {
            constructor(opts) {
                super(opts, "InflateRaw");
            }
        }
        class Unzip extends Zlib {
            constructor(opts) {
                super(opts, "Unzip");
            }
        }
        class Brotli extends ZlibBase {
            constructor(opts, mode) {
                opts = opts || {};
                opts.flush = opts.flush || constants.BROTLI_OPERATION_PROCESS;
                opts.finishFlush = opts.finishFlush || constants.BROTLI_OPERATION_FINISH;
                super(opts, mode);
                this[_fullFlushFlag] = constants.BROTLI_OPERATION_FLUSH;
            }
        }
        class BrotliCompress extends Brotli {
            constructor(opts) {
                super(opts, "BrotliCompress");
            }
        }
        class BrotliDecompress extends Brotli {
            constructor(opts) {
                super(opts, "BrotliDecompress");
            }
        }
        exports.Deflate = Deflate;
        exports.Inflate = Inflate;
        exports.Gzip = Gzip;
        exports.Gunzip = Gunzip;
        exports.DeflateRaw = DeflateRaw;
        exports.InflateRaw = InflateRaw;
        exports.Unzip = Unzip;
        if (typeof realZlib.BrotliCompress === "function") {
            exports.BrotliCompress = BrotliCompress;
            exports.BrotliDecompress = BrotliDecompress;
        } else {
            exports.BrotliCompress = exports.BrotliDecompress = class {
                constructor() {
                    throw new Error("Brotli is not supported in this version of Node.js");
                }
            };
        }
    },
    3179: (module, __unused_webpack_exports, __webpack_require__) => {
        const optsArg = __webpack_require__(2425);
        const pathArg = __webpack_require__(7394);
        const {mkdirpNative, mkdirpNativeSync} = __webpack_require__(5702);
        const {mkdirpManual, mkdirpManualSync} = __webpack_require__(8116);
        const {useNative, useNativeSync} = __webpack_require__(6631);
        const mkdirp = (path, opts) => {
            path = pathArg(path);
            opts = optsArg(opts);
            return useNative(opts) ? mkdirpNative(path, opts) : mkdirpManual(path, opts);
        };
        const mkdirpSync = (path, opts) => {
            path = pathArg(path);
            opts = optsArg(opts);
            return useNativeSync(opts) ? mkdirpNativeSync(path, opts) : mkdirpManualSync(path, opts);
        };
        mkdirp.sync = mkdirpSync;
        mkdirp.native = (path, opts) => mkdirpNative(pathArg(path), optsArg(opts));
        mkdirp.manual = (path, opts) => mkdirpManual(pathArg(path), optsArg(opts));
        mkdirp.nativeSync = (path, opts) => mkdirpNativeSync(pathArg(path), optsArg(opts));
        mkdirp.manualSync = (path, opts) => mkdirpManualSync(pathArg(path), optsArg(opts));
        module.exports = mkdirp;
    },
    1008: (module, __unused_webpack_exports, __webpack_require__) => {
        const {dirname} = __webpack_require__(5622);
        const findMade = (opts, parent, path = undefined) => {
            if (path === parent) return Promise.resolve();
            return opts.statAsync(parent).then((st => st.isDirectory() ? path : undefined), (er => er.code === "ENOENT" ? findMade(opts, dirname(parent), parent) : undefined));
        };
        const findMadeSync = (opts, parent, path = undefined) => {
            if (path === parent) return undefined;
            try {
                return opts.statSync(parent).isDirectory() ? path : undefined;
            } catch (er) {
                return er.code === "ENOENT" ? findMadeSync(opts, dirname(parent), parent) : undefined;
            }
        };
        module.exports = {
            findMade,
            findMadeSync
        };
    },
    8116: (module, __unused_webpack_exports, __webpack_require__) => {
        const {dirname} = __webpack_require__(5622);
        const mkdirpManual = (path, opts, made) => {
            opts.recursive = false;
            const parent = dirname(path);
            if (parent === path) {
                return opts.mkdirAsync(path, opts).catch((er => {
                    if (er.code !== "EISDIR") throw er;
                }));
            }
            return opts.mkdirAsync(path, opts).then((() => made || path), (er => {
                if (er.code === "ENOENT") return mkdirpManual(parent, opts).then((made => mkdirpManual(path, opts, made)));
                if (er.code !== "EEXIST" && er.code !== "EROFS") throw er;
                return opts.statAsync(path).then((st => {
                    if (st.isDirectory()) return made; else throw er;
                }), (() => {
                    throw er;
                }));
            }));
        };
        const mkdirpManualSync = (path, opts, made) => {
            const parent = dirname(path);
            opts.recursive = false;
            if (parent === path) {
                try {
                    return opts.mkdirSync(path, opts);
                } catch (er) {
                    if (er.code !== "EISDIR") throw er; else return;
                }
            }
            try {
                opts.mkdirSync(path, opts);
                return made || path;
            } catch (er) {
                if (er.code === "ENOENT") return mkdirpManualSync(path, opts, mkdirpManualSync(parent, opts, made));
                if (er.code !== "EEXIST" && er.code !== "EROFS") throw er;
                try {
                    if (!opts.statSync(path).isDirectory()) throw er;
                } catch (_) {
                    throw er;
                }
            }
        };
        module.exports = {
            mkdirpManual,
            mkdirpManualSync
        };
    },
    5702: (module, __unused_webpack_exports, __webpack_require__) => {
        const {dirname} = __webpack_require__(5622);
        const {findMade, findMadeSync} = __webpack_require__(1008);
        const {mkdirpManual, mkdirpManualSync} = __webpack_require__(8116);
        const mkdirpNative = (path, opts) => {
            opts.recursive = true;
            const parent = dirname(path);
            if (parent === path) return opts.mkdirAsync(path, opts);
            return findMade(opts, path).then((made => opts.mkdirAsync(path, opts).then((() => made)).catch((er => {
                if (er.code === "ENOENT") return mkdirpManual(path, opts); else throw er;
            }))));
        };
        const mkdirpNativeSync = (path, opts) => {
            opts.recursive = true;
            const parent = dirname(path);
            if (parent === path) return opts.mkdirSync(path, opts);
            const made = findMadeSync(opts, path);
            try {
                opts.mkdirSync(path, opts);
                return made;
            } catch (er) {
                if (er.code === "ENOENT") return mkdirpManualSync(path, opts); else throw er;
            }
        };
        module.exports = {
            mkdirpNative,
            mkdirpNativeSync
        };
    },
    2425: (module, __unused_webpack_exports, __webpack_require__) => {
        const {promisify} = __webpack_require__(1669);
        const fs = __webpack_require__(5747);
        const optsArg = opts => {
            if (!opts) opts = {
                mode: 511,
                fs
            }; else if (typeof opts === "object") opts = {
                mode: 511,
                fs,
                ...opts
            }; else if (typeof opts === "number") opts = {
                mode: opts,
                fs
            }; else if (typeof opts === "string") opts = {
                mode: parseInt(opts, 8),
                fs
            }; else throw new TypeError("invalid options argument");
            opts.mkdir = opts.mkdir || opts.fs.mkdir || fs.mkdir;
            opts.mkdirAsync = promisify(opts.mkdir);
            opts.stat = opts.stat || opts.fs.stat || fs.stat;
            opts.statAsync = promisify(opts.stat);
            opts.statSync = opts.statSync || opts.fs.statSync || fs.statSync;
            opts.mkdirSync = opts.mkdirSync || opts.fs.mkdirSync || fs.mkdirSync;
            return opts;
        };
        module.exports = optsArg;
    },
    7394: (module, __unused_webpack_exports, __webpack_require__) => {
        const platform = process.env.__TESTING_MKDIRP_PLATFORM__ || process.platform;
        const {resolve, parse} = __webpack_require__(5622);
        const pathArg = path => {
            if (/\0/.test(path)) {
                throw Object.assign(new TypeError("path must be a string without null bytes"), {
                    path,
                    code: "ERR_INVALID_ARG_VALUE"
                });
            }
            path = resolve(path);
            if (platform === "win32") {
                const badWinChars = /[*|"<>?:]/;
                const {root} = parse(path);
                if (badWinChars.test(path.substr(root.length))) {
                    throw Object.assign(new Error("Illegal characters in path."), {
                        path,
                        code: "EINVAL"
                    });
                }
            }
            return path;
        };
        module.exports = pathArg;
    },
    6631: (module, __unused_webpack_exports, __webpack_require__) => {
        const fs = __webpack_require__(5747);
        const version = process.env.__TESTING_MKDIRP_NODE_VERSION__ || process.version;
        const versArr = version.replace(/^v/, "").split(".");
        const hasNative = +versArr[0] > 10 || +versArr[0] === 10 && +versArr[1] >= 12;
        const useNative = !hasNative ? () => false : opts => opts.mkdir === fs.mkdir;
        const useNativeSync = !hasNative ? () => false : opts => opts.mkdirSync === fs.mkdirSync;
        module.exports = {
            useNative,
            useNativeSync
        };
    },
    1189: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        exports.c = exports.create = __webpack_require__(9540);
        exports.r = exports.replace = __webpack_require__(3666);
        exports.t = exports.list = __webpack_require__(1090);
        exports.u = exports.update = __webpack_require__(4229);
        exports.x = exports.extract = __webpack_require__(1372);
        exports.Pack = __webpack_require__(5843);
        exports.Unpack = __webpack_require__(2864);
        exports.Parse = __webpack_require__(6234);
        exports.ReadEntry = __webpack_require__(7847);
        exports.WriteEntry = __webpack_require__(8418);
        exports.Header = __webpack_require__(5017);
        exports.Pax = __webpack_require__(9154);
        exports.types = __webpack_require__(9806);
    },
    9540: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const hlo = __webpack_require__(7461);
        const Pack = __webpack_require__(5843);
        const fsm = __webpack_require__(8553);
        const t = __webpack_require__(1090);
        const path = __webpack_require__(5622);
        module.exports = (opt_, files, cb) => {
            if (typeof files === "function") cb = files;
            if (Array.isArray(opt_)) files = opt_, opt_ = {};
            if (!files || !Array.isArray(files) || !files.length) throw new TypeError("no files or directories specified");
            files = Array.from(files);
            const opt = hlo(opt_);
            if (opt.sync && typeof cb === "function") throw new TypeError("callback not supported for sync tar functions");
            if (!opt.file && typeof cb === "function") throw new TypeError("callback only supported with file option");
            return opt.file && opt.sync ? createFileSync(opt, files) : opt.file ? createFile(opt, files, cb) : opt.sync ? createSync(opt, files) : create(opt, files);
        };
        const createFileSync = (opt, files) => {
            const p = new Pack.Sync(opt);
            const stream = new fsm.WriteStreamSync(opt.file, {
                mode: opt.mode || 438
            });
            p.pipe(stream);
            addFilesSync(p, files);
        };
        const createFile = (opt, files, cb) => {
            const p = new Pack(opt);
            const stream = new fsm.WriteStream(opt.file, {
                mode: opt.mode || 438
            });
            p.pipe(stream);
            const promise = new Promise(((res, rej) => {
                stream.on("error", rej);
                stream.on("close", res);
                p.on("error", rej);
            }));
            addFilesAsync(p, files);
            return cb ? promise.then(cb, cb) : promise;
        };
        const addFilesSync = (p, files) => {
            files.forEach((file => {
                if (file.charAt(0) === "@") {
                    t({
                        file: path.resolve(p.cwd, file.substr(1)),
                        sync: true,
                        noResume: true,
                        onentry: entry => p.add(entry)
                    });
                } else p.add(file);
            }));
            p.end();
        };
        const addFilesAsync = (p, files) => {
            while (files.length) {
                const file = files.shift();
                if (file.charAt(0) === "@") {
                    return t({
                        file: path.resolve(p.cwd, file.substr(1)),
                        noResume: true,
                        onentry: entry => p.add(entry)
                    }).then((_ => addFilesAsync(p, files)));
                } else p.add(file);
            }
            p.end();
        };
        const createSync = (opt, files) => {
            const p = new Pack.Sync(opt);
            addFilesSync(p, files);
            return p;
        };
        const create = (opt, files) => {
            const p = new Pack(opt);
            addFilesAsync(p, files);
            return p;
        };
    },
    1372: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const hlo = __webpack_require__(7461);
        const Unpack = __webpack_require__(2864);
        const fs = __webpack_require__(5747);
        const fsm = __webpack_require__(8553);
        const path = __webpack_require__(5622);
        module.exports = (opt_, files, cb) => {
            if (typeof opt_ === "function") cb = opt_, files = null, opt_ = {}; else if (Array.isArray(opt_)) files = opt_, 
            opt_ = {};
            if (typeof files === "function") cb = files, files = null;
            if (!files) files = []; else files = Array.from(files);
            const opt = hlo(opt_);
            if (opt.sync && typeof cb === "function") throw new TypeError("callback not supported for sync tar functions");
            if (!opt.file && typeof cb === "function") throw new TypeError("callback only supported with file option");
            if (files.length) filesFilter(opt, files);
            return opt.file && opt.sync ? extractFileSync(opt) : opt.file ? extractFile(opt, cb) : opt.sync ? extractSync(opt) : extract(opt);
        };
        const filesFilter = (opt, files) => {
            const map = new Map(files.map((f => [ f.replace(/\/+$/, ""), true ])));
            const filter = opt.filter;
            const mapHas = (file, r) => {
                const root = r || path.parse(file).root || ".";
                const ret = file === root ? false : map.has(file) ? map.get(file) : mapHas(path.dirname(file), root);
                map.set(file, ret);
                return ret;
            };
            opt.filter = filter ? (file, entry) => filter(file, entry) && mapHas(file.replace(/\/+$/, "")) : file => mapHas(file.replace(/\/+$/, ""));
        };
        const extractFileSync = opt => {
            const u = new Unpack.Sync(opt);
            const file = opt.file;
            const stat = fs.statSync(file);
            const readSize = opt.maxReadSize || 16 * 1024 * 1024;
            const stream = new fsm.ReadStreamSync(file, {
                readSize,
                size: stat.size
            });
            stream.pipe(u);
        };
        const extractFile = (opt, cb) => {
            const u = new Unpack(opt);
            const readSize = opt.maxReadSize || 16 * 1024 * 1024;
            const file = opt.file;
            const p = new Promise(((resolve, reject) => {
                u.on("error", reject);
                u.on("close", resolve);
                fs.stat(file, ((er, stat) => {
                    if (er) reject(er); else {
                        const stream = new fsm.ReadStream(file, {
                            readSize,
                            size: stat.size
                        });
                        stream.on("error", reject);
                        stream.pipe(u);
                    }
                }));
            }));
            return cb ? p.then(cb, cb) : p;
        };
        const extractSync = opt => new Unpack.Sync(opt);
        const extract = opt => new Unpack(opt);
    },
    8512: (module, __unused_webpack_exports, __webpack_require__) => {
        const platform = process.env.__FAKE_PLATFORM__ || process.platform;
        const isWindows = platform === "win32";
        const fs = global.__FAKE_TESTING_FS__ || __webpack_require__(5747);
        const {O_CREAT, O_TRUNC, O_WRONLY, UV_FS_O_FILEMAP = 0} = fs.constants;
        const fMapEnabled = isWindows && !!UV_FS_O_FILEMAP;
        const fMapLimit = 512 * 1024;
        const fMapFlag = UV_FS_O_FILEMAP | O_TRUNC | O_CREAT | O_WRONLY;
        module.exports = !fMapEnabled ? () => "w" : size => size < fMapLimit ? fMapFlag : "w";
    },
    5017: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const types = __webpack_require__(9806);
        const pathModule = __webpack_require__(5622).posix;
        const large = __webpack_require__(2795);
        const SLURP = Symbol("slurp");
        const TYPE = Symbol("type");
        class Header {
            constructor(data, off, ex, gex) {
                this.cksumValid = false;
                this.needPax = false;
                this.nullBlock = false;
                this.block = null;
                this.path = null;
                this.mode = null;
                this.uid = null;
                this.gid = null;
                this.size = null;
                this.mtime = null;
                this.cksum = null;
                this[TYPE] = "0";
                this.linkpath = null;
                this.uname = null;
                this.gname = null;
                this.devmaj = 0;
                this.devmin = 0;
                this.atime = null;
                this.ctime = null;
                if (Buffer.isBuffer(data)) this.decode(data, off || 0, ex, gex); else if (data) this.set(data);
            }
            decode(buf, off, ex, gex) {
                if (!off) off = 0;
                if (!buf || !(buf.length >= off + 512)) throw new Error("need 512 bytes for header");
                this.path = decString(buf, off, 100);
                this.mode = decNumber(buf, off + 100, 8);
                this.uid = decNumber(buf, off + 108, 8);
                this.gid = decNumber(buf, off + 116, 8);
                this.size = decNumber(buf, off + 124, 12);
                this.mtime = decDate(buf, off + 136, 12);
                this.cksum = decNumber(buf, off + 148, 12);
                this[SLURP](ex);
                this[SLURP](gex, true);
                this[TYPE] = decString(buf, off + 156, 1);
                if (this[TYPE] === "") this[TYPE] = "0";
                if (this[TYPE] === "0" && this.path.substr(-1) === "/") this[TYPE] = "5";
                if (this[TYPE] === "5") this.size = 0;
                this.linkpath = decString(buf, off + 157, 100);
                if (buf.slice(off + 257, off + 265).toString() === "ustar\x0000") {
                    this.uname = decString(buf, off + 265, 32);
                    this.gname = decString(buf, off + 297, 32);
                    this.devmaj = decNumber(buf, off + 329, 8);
                    this.devmin = decNumber(buf, off + 337, 8);
                    if (buf[off + 475] !== 0) {
                        const prefix = decString(buf, off + 345, 155);
                        this.path = prefix + "/" + this.path;
                    } else {
                        const prefix = decString(buf, off + 345, 130);
                        if (prefix) this.path = prefix + "/" + this.path;
                        this.atime = decDate(buf, off + 476, 12);
                        this.ctime = decDate(buf, off + 488, 12);
                    }
                }
                let sum = 8 * 32;
                for (let i = off; i < off + 148; i++) sum += buf[i];
                for (let i = off + 156; i < off + 512; i++) sum += buf[i];
                this.cksumValid = sum === this.cksum;
                if (this.cksum === null && sum === 8 * 32) this.nullBlock = true;
            }
            [SLURP](ex, global) {
                for (const k in ex) {
                    if (ex[k] !== null && ex[k] !== undefined && !(global && k === "path")) this[k] = ex[k];
                }
            }
            encode(buf, off) {
                if (!buf) {
                    buf = this.block = Buffer.alloc(512);
                    off = 0;
                }
                if (!off) off = 0;
                if (!(buf.length >= off + 512)) throw new Error("need 512 bytes for header");
                const prefixSize = this.ctime || this.atime ? 130 : 155;
                const split = splitPrefix(this.path || "", prefixSize);
                const path = split[0];
                const prefix = split[1];
                this.needPax = split[2];
                this.needPax = encString(buf, off, 100, path) || this.needPax;
                this.needPax = encNumber(buf, off + 100, 8, this.mode) || this.needPax;
                this.needPax = encNumber(buf, off + 108, 8, this.uid) || this.needPax;
                this.needPax = encNumber(buf, off + 116, 8, this.gid) || this.needPax;
                this.needPax = encNumber(buf, off + 124, 12, this.size) || this.needPax;
                this.needPax = encDate(buf, off + 136, 12, this.mtime) || this.needPax;
                buf[off + 156] = this[TYPE].charCodeAt(0);
                this.needPax = encString(buf, off + 157, 100, this.linkpath) || this.needPax;
                buf.write("ustar\x0000", off + 257, 8);
                this.needPax = encString(buf, off + 265, 32, this.uname) || this.needPax;
                this.needPax = encString(buf, off + 297, 32, this.gname) || this.needPax;
                this.needPax = encNumber(buf, off + 329, 8, this.devmaj) || this.needPax;
                this.needPax = encNumber(buf, off + 337, 8, this.devmin) || this.needPax;
                this.needPax = encString(buf, off + 345, prefixSize, prefix) || this.needPax;
                if (buf[off + 475] !== 0) this.needPax = encString(buf, off + 345, 155, prefix) || this.needPax; else {
                    this.needPax = encString(buf, off + 345, 130, prefix) || this.needPax;
                    this.needPax = encDate(buf, off + 476, 12, this.atime) || this.needPax;
                    this.needPax = encDate(buf, off + 488, 12, this.ctime) || this.needPax;
                }
                let sum = 8 * 32;
                for (let i = off; i < off + 148; i++) sum += buf[i];
                for (let i = off + 156; i < off + 512; i++) sum += buf[i];
                this.cksum = sum;
                encNumber(buf, off + 148, 8, this.cksum);
                this.cksumValid = true;
                return this.needPax;
            }
            set(data) {
                for (const i in data) {
                    if (data[i] !== null && data[i] !== undefined) this[i] = data[i];
                }
            }
            get type() {
                return types.name.get(this[TYPE]) || this[TYPE];
            }
            get typeKey() {
                return this[TYPE];
            }
            set type(type) {
                if (types.code.has(type)) this[TYPE] = types.code.get(type); else this[TYPE] = type;
            }
        }
        const splitPrefix = (p, prefixSize) => {
            const pathSize = 100;
            let pp = p;
            let prefix = "";
            let ret;
            const root = pathModule.parse(p).root || ".";
            if (Buffer.byteLength(pp) < pathSize) ret = [ pp, prefix, false ]; else {
                prefix = pathModule.dirname(pp);
                pp = pathModule.basename(pp);
                do {
                    if (Buffer.byteLength(pp) <= pathSize && Buffer.byteLength(prefix) <= prefixSize) ret = [ pp, prefix, false ]; else if (Buffer.byteLength(pp) > pathSize && Buffer.byteLength(prefix) <= prefixSize) ret = [ pp.substr(0, pathSize - 1), prefix, true ]; else {
                        pp = pathModule.join(pathModule.basename(prefix), pp);
                        prefix = pathModule.dirname(prefix);
                    }
                } while (prefix !== root && !ret);
                if (!ret) ret = [ p.substr(0, pathSize - 1), "", true ];
            }
            return ret;
        };
        const decString = (buf, off, size) => buf.slice(off, off + size).toString("utf8").replace(/\0.*/, "");
        const decDate = (buf, off, size) => numToDate(decNumber(buf, off, size));
        const numToDate = num => num === null ? null : new Date(num * 1e3);
        const decNumber = (buf, off, size) => buf[off] & 128 ? large.parse(buf.slice(off, off + size)) : decSmallNumber(buf, off, size);
        const nanNull = value => isNaN(value) ? null : value;
        const decSmallNumber = (buf, off, size) => nanNull(parseInt(buf.slice(off, off + size).toString("utf8").replace(/\0.*$/, "").trim(), 8));
        const MAXNUM = {
            12: 8589934591,
            8: 2097151
        };
        const encNumber = (buf, off, size, number) => number === null ? false : number > MAXNUM[size] || number < 0 ? (large.encode(number, buf.slice(off, off + size)), 
        true) : (encSmallNumber(buf, off, size, number), false);
        const encSmallNumber = (buf, off, size, number) => buf.write(octalString(number, size), off, size, "ascii");
        const octalString = (number, size) => padOctal(Math.floor(number).toString(8), size);
        const padOctal = (string, size) => (string.length === size - 1 ? string : new Array(size - string.length - 1).join("0") + string + " ") + "\0";
        const encDate = (buf, off, size, date) => date === null ? false : encNumber(buf, off, size, date.getTime() / 1e3);
        const NULLS = new Array(156).join("\0");
        const encString = (buf, off, size, string) => string === null ? false : (buf.write(string + NULLS, off, size, "utf8"), 
        string.length !== Buffer.byteLength(string) || string.length > size);
        module.exports = Header;
    },
    7461: module => {
        "use strict";
        const argmap = new Map([ [ "C", "cwd" ], [ "f", "file" ], [ "z", "gzip" ], [ "P", "preservePaths" ], [ "U", "unlink" ], [ "strip-components", "strip" ], [ "stripComponents", "strip" ], [ "keep-newer", "newer" ], [ "keepNewer", "newer" ], [ "keep-newer-files", "newer" ], [ "keepNewerFiles", "newer" ], [ "k", "keep" ], [ "keep-existing", "keep" ], [ "keepExisting", "keep" ], [ "m", "noMtime" ], [ "no-mtime", "noMtime" ], [ "p", "preserveOwner" ], [ "L", "follow" ], [ "h", "follow" ] ]);
        module.exports = opt => opt ? Object.keys(opt).map((k => [ argmap.has(k) ? argmap.get(k) : k, opt[k] ])).reduce(((set, kv) => (set[kv[0]] = kv[1], 
        set)), Object.create(null)) : {};
    },
    2795: module => {
        "use strict";
        const encode = (num, buf) => {
            if (!Number.isSafeInteger(num)) throw Error("cannot encode number outside of javascript safe integer range"); else if (num < 0) encodeNegative(num, buf); else encodePositive(num, buf);
            return buf;
        };
        const encodePositive = (num, buf) => {
            buf[0] = 128;
            for (var i = buf.length; i > 1; i--) {
                buf[i - 1] = num & 255;
                num = Math.floor(num / 256);
            }
        };
        const encodeNegative = (num, buf) => {
            buf[0] = 255;
            var flipped = false;
            num = num * -1;
            for (var i = buf.length; i > 1; i--) {
                var byte = num & 255;
                num = Math.floor(num / 256);
                if (flipped) buf[i - 1] = onesComp(byte); else if (byte === 0) buf[i - 1] = 0; else {
                    flipped = true;
                    buf[i - 1] = twosComp(byte);
                }
            }
        };
        const parse = buf => {
            const pre = buf[0];
            const value = pre === 128 ? pos(buf.slice(1, buf.length)) : pre === 255 ? twos(buf) : null;
            if (value === null) throw Error("invalid base256 encoding");
            if (!Number.isSafeInteger(value)) throw Error("parsed number outside of javascript safe integer range");
            return value;
        };
        const twos = buf => {
            var len = buf.length;
            var sum = 0;
            var flipped = false;
            for (var i = len - 1; i > -1; i--) {
                var byte = buf[i];
                var f;
                if (flipped) f = onesComp(byte); else if (byte === 0) f = byte; else {
                    flipped = true;
                    f = twosComp(byte);
                }
                if (f !== 0) sum -= f * Math.pow(256, len - i - 1);
            }
            return sum;
        };
        const pos = buf => {
            var len = buf.length;
            var sum = 0;
            for (var i = len - 1; i > -1; i--) {
                var byte = buf[i];
                if (byte !== 0) sum += byte * Math.pow(256, len - i - 1);
            }
            return sum;
        };
        const onesComp = byte => (255 ^ byte) & 255;
        const twosComp = byte => (255 ^ byte) + 1 & 255;
        module.exports = {
            encode,
            parse
        };
    },
    1090: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const hlo = __webpack_require__(7461);
        const Parser = __webpack_require__(6234);
        const fs = __webpack_require__(5747);
        const fsm = __webpack_require__(8553);
        const path = __webpack_require__(5622);
        module.exports = (opt_, files, cb) => {
            if (typeof opt_ === "function") cb = opt_, files = null, opt_ = {}; else if (Array.isArray(opt_)) files = opt_, 
            opt_ = {};
            if (typeof files === "function") cb = files, files = null;
            if (!files) files = []; else files = Array.from(files);
            const opt = hlo(opt_);
            if (opt.sync && typeof cb === "function") throw new TypeError("callback not supported for sync tar functions");
            if (!opt.file && typeof cb === "function") throw new TypeError("callback only supported with file option");
            if (files.length) filesFilter(opt, files);
            if (!opt.noResume) onentryFunction(opt);
            return opt.file && opt.sync ? listFileSync(opt) : opt.file ? listFile(opt, cb) : list(opt);
        };
        const onentryFunction = opt => {
            const onentry = opt.onentry;
            opt.onentry = onentry ? e => {
                onentry(e);
                e.resume();
            } : e => e.resume();
        };
        const filesFilter = (opt, files) => {
            const map = new Map(files.map((f => [ f.replace(/\/+$/, ""), true ])));
            const filter = opt.filter;
            const mapHas = (file, r) => {
                const root = r || path.parse(file).root || ".";
                const ret = file === root ? false : map.has(file) ? map.get(file) : mapHas(path.dirname(file), root);
                map.set(file, ret);
                return ret;
            };
            opt.filter = filter ? (file, entry) => filter(file, entry) && mapHas(file.replace(/\/+$/, "")) : file => mapHas(file.replace(/\/+$/, ""));
        };
        const listFileSync = opt => {
            const p = list(opt);
            const file = opt.file;
            let threw = true;
            let fd;
            try {
                const stat = fs.statSync(file);
                const readSize = opt.maxReadSize || 16 * 1024 * 1024;
                if (stat.size < readSize) p.end(fs.readFileSync(file)); else {
                    let pos = 0;
                    const buf = Buffer.allocUnsafe(readSize);
                    fd = fs.openSync(file, "r");
                    while (pos < stat.size) {
                        const bytesRead = fs.readSync(fd, buf, 0, readSize, pos);
                        pos += bytesRead;
                        p.write(buf.slice(0, bytesRead));
                    }
                    p.end();
                }
                threw = false;
            } finally {
                if (threw && fd) {
                    try {
                        fs.closeSync(fd);
                    } catch (er) {}
                }
            }
        };
        const listFile = (opt, cb) => {
            const parse = new Parser(opt);
            const readSize = opt.maxReadSize || 16 * 1024 * 1024;
            const file = opt.file;
            const p = new Promise(((resolve, reject) => {
                parse.on("error", reject);
                parse.on("end", resolve);
                fs.stat(file, ((er, stat) => {
                    if (er) reject(er); else {
                        const stream = new fsm.ReadStream(file, {
                            readSize,
                            size: stat.size
                        });
                        stream.on("error", reject);
                        stream.pipe(parse);
                    }
                }));
            }));
            return cb ? p.then(cb, cb) : p;
        };
        const list = opt => new Parser(opt);
    },
    3956: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const mkdirp = __webpack_require__(3179);
        const fs = __webpack_require__(5747);
        const path = __webpack_require__(5622);
        const chownr = __webpack_require__(2047);
        class SymlinkError extends Error {
            constructor(symlink, path) {
                super("Cannot extract through symbolic link");
                this.path = path;
                this.symlink = symlink;
            }
            get name() {
                return "SylinkError";
            }
        }
        class CwdError extends Error {
            constructor(path, code) {
                super(code + ": Cannot cd into '" + path + "'");
                this.path = path;
                this.code = code;
            }
            get name() {
                return "CwdError";
            }
        }
        module.exports = (dir, opt, cb) => {
            const umask = opt.umask;
            const mode = opt.mode | 448;
            const needChmod = (mode & umask) !== 0;
            const uid = opt.uid;
            const gid = opt.gid;
            const doChown = typeof uid === "number" && typeof gid === "number" && (uid !== opt.processUid || gid !== opt.processGid);
            const preserve = opt.preserve;
            const unlink = opt.unlink;
            const cache = opt.cache;
            const cwd = opt.cwd;
            const done = (er, created) => {
                if (er) cb(er); else {
                    cache.set(dir, true);
                    if (created && doChown) chownr(created, uid, gid, (er => done(er))); else if (needChmod) fs.chmod(dir, mode, cb); else cb();
                }
            };
            if (cache && cache.get(dir) === true) return done();
            if (dir === cwd) {
                return fs.stat(dir, ((er, st) => {
                    if (er || !st.isDirectory()) er = new CwdError(dir, er && er.code || "ENOTDIR");
                    done(er);
                }));
            }
            if (preserve) return mkdirp(dir, {
                mode
            }).then((made => done(null, made)), done);
            const sub = path.relative(cwd, dir);
            const parts = sub.split(/\/|\\/);
            mkdir_(cwd, parts, mode, cache, unlink, cwd, null, done);
        };
        const mkdir_ = (base, parts, mode, cache, unlink, cwd, created, cb) => {
            if (!parts.length) return cb(null, created);
            const p = parts.shift();
            const part = base + "/" + p;
            if (cache.get(part)) return mkdir_(part, parts, mode, cache, unlink, cwd, created, cb);
            fs.mkdir(part, mode, onmkdir(part, parts, mode, cache, unlink, cwd, created, cb));
        };
        const onmkdir = (part, parts, mode, cache, unlink, cwd, created, cb) => er => {
            if (er) {
                if (er.path && path.dirname(er.path) === cwd && (er.code === "ENOTDIR" || er.code === "ENOENT")) return cb(new CwdError(cwd, er.code));
                fs.lstat(part, ((statEr, st) => {
                    if (statEr) cb(statEr); else if (st.isDirectory()) mkdir_(part, parts, mode, cache, unlink, cwd, created, cb); else if (unlink) {
                        fs.unlink(part, (er => {
                            if (er) return cb(er);
                            fs.mkdir(part, mode, onmkdir(part, parts, mode, cache, unlink, cwd, created, cb));
                        }));
                    } else if (st.isSymbolicLink()) return cb(new SymlinkError(part, part + "/" + parts.join("/"))); else cb(er);
                }));
            } else {
                created = created || part;
                mkdir_(part, parts, mode, cache, unlink, cwd, created, cb);
            }
        };
        module.exports.sync = (dir, opt) => {
            const umask = opt.umask;
            const mode = opt.mode | 448;
            const needChmod = (mode & umask) !== 0;
            const uid = opt.uid;
            const gid = opt.gid;
            const doChown = typeof uid === "number" && typeof gid === "number" && (uid !== opt.processUid || gid !== opt.processGid);
            const preserve = opt.preserve;
            const unlink = opt.unlink;
            const cache = opt.cache;
            const cwd = opt.cwd;
            const done = created => {
                cache.set(dir, true);
                if (created && doChown) chownr.sync(created, uid, gid);
                if (needChmod) fs.chmodSync(dir, mode);
            };
            if (cache && cache.get(dir) === true) return done();
            if (dir === cwd) {
                let ok = false;
                let code = "ENOTDIR";
                try {
                    ok = fs.statSync(dir).isDirectory();
                } catch (er) {
                    code = er.code;
                } finally {
                    if (!ok) throw new CwdError(dir, code);
                }
                done();
                return;
            }
            if (preserve) return done(mkdirp.sync(dir, mode));
            const sub = path.relative(cwd, dir);
            const parts = sub.split(/\/|\\/);
            let created = null;
            for (let p = parts.shift(), part = cwd; p && (part += "/" + p); p = parts.shift()) {
                if (cache.get(part)) continue;
                try {
                    fs.mkdirSync(part, mode);
                    created = created || part;
                    cache.set(part, true);
                } catch (er) {
                    if (er.path && path.dirname(er.path) === cwd && (er.code === "ENOTDIR" || er.code === "ENOENT")) return new CwdError(cwd, er.code);
                    const st = fs.lstatSync(part);
                    if (st.isDirectory()) {
                        cache.set(part, true);
                        continue;
                    } else if (unlink) {
                        fs.unlinkSync(part);
                        fs.mkdirSync(part, mode);
                        created = created || part;
                        cache.set(part, true);
                        continue;
                    } else if (st.isSymbolicLink()) return new SymlinkError(part, part + "/" + parts.join("/"));
                }
            }
            return done(created);
        };
    },
    9574: module => {
        "use strict";
        module.exports = (mode, isDir, portable) => {
            mode &= 4095;
            if (portable) mode = (mode | 384) & ~18;
            if (isDir) {
                if (mode & 256) mode |= 64;
                if (mode & 32) mode |= 8;
                if (mode & 4) mode |= 1;
            }
            return mode;
        };
    },
    5843: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        class PackJob {
            constructor(path, absolute) {
                this.path = path || "./";
                this.absolute = absolute;
                this.entry = null;
                this.stat = null;
                this.readdir = null;
                this.pending = false;
                this.ignore = false;
                this.piped = false;
            }
        }
        const MiniPass = __webpack_require__(2253);
        const zlib = __webpack_require__(3704);
        const ReadEntry = __webpack_require__(7847);
        const WriteEntry = __webpack_require__(8418);
        const WriteEntrySync = WriteEntry.Sync;
        const WriteEntryTar = WriteEntry.Tar;
        const Yallist = __webpack_require__(1455);
        const EOF = Buffer.alloc(1024);
        const ONSTAT = Symbol("onStat");
        const ENDED = Symbol("ended");
        const QUEUE = Symbol("queue");
        const CURRENT = Symbol("current");
        const PROCESS = Symbol("process");
        const PROCESSING = Symbol("processing");
        const PROCESSJOB = Symbol("processJob");
        const JOBS = Symbol("jobs");
        const JOBDONE = Symbol("jobDone");
        const ADDFSENTRY = Symbol("addFSEntry");
        const ADDTARENTRY = Symbol("addTarEntry");
        const STAT = Symbol("stat");
        const READDIR = Symbol("readdir");
        const ONREADDIR = Symbol("onreaddir");
        const PIPE = Symbol("pipe");
        const ENTRY = Symbol("entry");
        const ENTRYOPT = Symbol("entryOpt");
        const WRITEENTRYCLASS = Symbol("writeEntryClass");
        const WRITE = Symbol("write");
        const ONDRAIN = Symbol("ondrain");
        const fs = __webpack_require__(5747);
        const path = __webpack_require__(5622);
        const warner = __webpack_require__(8783);
        const Pack = warner(class Pack extends MiniPass {
            constructor(opt) {
                super(opt);
                opt = opt || Object.create(null);
                this.opt = opt;
                this.file = opt.file || "";
                this.cwd = opt.cwd || process.cwd();
                this.maxReadSize = opt.maxReadSize;
                this.preservePaths = !!opt.preservePaths;
                this.strict = !!opt.strict;
                this.noPax = !!opt.noPax;
                this.prefix = (opt.prefix || "").replace(/(\\|\/)+$/, "");
                this.linkCache = opt.linkCache || new Map;
                this.statCache = opt.statCache || new Map;
                this.readdirCache = opt.readdirCache || new Map;
                this[WRITEENTRYCLASS] = WriteEntry;
                if (typeof opt.onwarn === "function") this.on("warn", opt.onwarn);
                this.portable = !!opt.portable;
                this.zip = null;
                if (opt.gzip) {
                    if (typeof opt.gzip !== "object") opt.gzip = {};
                    if (this.portable) opt.gzip.portable = true;
                    this.zip = new zlib.Gzip(opt.gzip);
                    this.zip.on("data", (chunk => super.write(chunk)));
                    this.zip.on("end", (_ => super.end()));
                    this.zip.on("drain", (_ => this[ONDRAIN]()));
                    this.on("resume", (_ => this.zip.resume()));
                } else this.on("drain", this[ONDRAIN]);
                this.noDirRecurse = !!opt.noDirRecurse;
                this.follow = !!opt.follow;
                this.noMtime = !!opt.noMtime;
                this.mtime = opt.mtime || null;
                this.filter = typeof opt.filter === "function" ? opt.filter : _ => true;
                this[QUEUE] = new Yallist;
                this[JOBS] = 0;
                this.jobs = +opt.jobs || 4;
                this[PROCESSING] = false;
                this[ENDED] = false;
            }
            [WRITE](chunk) {
                return super.write(chunk);
            }
            add(path) {
                this.write(path);
                return this;
            }
            end(path) {
                if (path) this.write(path);
                this[ENDED] = true;
                this[PROCESS]();
                return this;
            }
            write(path) {
                if (this[ENDED]) throw new Error("write after end");
                if (path instanceof ReadEntry) this[ADDTARENTRY](path); else this[ADDFSENTRY](path);
                return this.flowing;
            }
            [ADDTARENTRY](p) {
                const absolute = path.resolve(this.cwd, p.path);
                if (this.prefix) p.path = this.prefix + "/" + p.path.replace(/^\.(\/+|$)/, "");
                if (!this.filter(p.path, p)) p.resume(); else {
                    const job = new PackJob(p.path, absolute, false);
                    job.entry = new WriteEntryTar(p, this[ENTRYOPT](job));
                    job.entry.on("end", (_ => this[JOBDONE](job)));
                    this[JOBS] += 1;
                    this[QUEUE].push(job);
                }
                this[PROCESS]();
            }
            [ADDFSENTRY](p) {
                const absolute = path.resolve(this.cwd, p);
                if (this.prefix) p = this.prefix + "/" + p.replace(/^\.(\/+|$)/, "");
                this[QUEUE].push(new PackJob(p, absolute));
                this[PROCESS]();
            }
            [STAT](job) {
                job.pending = true;
                this[JOBS] += 1;
                const stat = this.follow ? "stat" : "lstat";
                fs[stat](job.absolute, ((er, stat) => {
                    job.pending = false;
                    this[JOBS] -= 1;
                    if (er) this.emit("error", er); else this[ONSTAT](job, stat);
                }));
            }
            [ONSTAT](job, stat) {
                this.statCache.set(job.absolute, stat);
                job.stat = stat;
                if (!this.filter(job.path, stat)) job.ignore = true;
                this[PROCESS]();
            }
            [READDIR](job) {
                job.pending = true;
                this[JOBS] += 1;
                fs.readdir(job.absolute, ((er, entries) => {
                    job.pending = false;
                    this[JOBS] -= 1;
                    if (er) return this.emit("error", er);
                    this[ONREADDIR](job, entries);
                }));
            }
            [ONREADDIR](job, entries) {
                this.readdirCache.set(job.absolute, entries);
                job.readdir = entries;
                this[PROCESS]();
            }
            [PROCESS]() {
                if (this[PROCESSING]) return;
                this[PROCESSING] = true;
                for (let w = this[QUEUE].head; w !== null && this[JOBS] < this.jobs; w = w.next) {
                    this[PROCESSJOB](w.value);
                    if (w.value.ignore) {
                        const p = w.next;
                        this[QUEUE].removeNode(w);
                        w.next = p;
                    }
                }
                this[PROCESSING] = false;
                if (this[ENDED] && !this[QUEUE].length && this[JOBS] === 0) {
                    if (this.zip) this.zip.end(EOF); else {
                        super.write(EOF);
                        super.end();
                    }
                }
            }
            get [CURRENT]() {
                return this[QUEUE] && this[QUEUE].head && this[QUEUE].head.value;
            }
            [JOBDONE](job) {
                this[QUEUE].shift();
                this[JOBS] -= 1;
                this[PROCESS]();
            }
            [PROCESSJOB](job) {
                if (job.pending) return;
                if (job.entry) {
                    if (job === this[CURRENT] && !job.piped) this[PIPE](job);
                    return;
                }
                if (!job.stat) {
                    if (this.statCache.has(job.absolute)) this[ONSTAT](job, this.statCache.get(job.absolute)); else this[STAT](job);
                }
                if (!job.stat) return;
                if (job.ignore) return;
                if (!this.noDirRecurse && job.stat.isDirectory() && !job.readdir) {
                    if (this.readdirCache.has(job.absolute)) this[ONREADDIR](job, this.readdirCache.get(job.absolute)); else this[READDIR](job);
                    if (!job.readdir) return;
                }
                job.entry = this[ENTRY](job);
                if (!job.entry) {
                    job.ignore = true;
                    return;
                }
                if (job === this[CURRENT] && !job.piped) this[PIPE](job);
            }
            [ENTRYOPT](job) {
                return {
                    onwarn: (code, msg, data) => this.warn(code, msg, data),
                    noPax: this.noPax,
                    cwd: this.cwd,
                    absolute: job.absolute,
                    preservePaths: this.preservePaths,
                    maxReadSize: this.maxReadSize,
                    strict: this.strict,
                    portable: this.portable,
                    linkCache: this.linkCache,
                    statCache: this.statCache,
                    noMtime: this.noMtime,
                    mtime: this.mtime
                };
            }
            [ENTRY](job) {
                this[JOBS] += 1;
                try {
                    return new this[WRITEENTRYCLASS](job.path, this[ENTRYOPT](job)).on("end", (() => this[JOBDONE](job))).on("error", (er => this.emit("error", er)));
                } catch (er) {
                    this.emit("error", er);
                }
            }
            [ONDRAIN]() {
                if (this[CURRENT] && this[CURRENT].entry) this[CURRENT].entry.resume();
            }
            [PIPE](job) {
                job.piped = true;
                if (job.readdir) {
                    job.readdir.forEach((entry => {
                        const p = this.prefix ? job.path.slice(this.prefix.length + 1) || "./" : job.path;
                        const base = p === "./" ? "" : p.replace(/\/*$/, "/");
                        this[ADDFSENTRY](base + entry);
                    }));
                }
                const source = job.entry;
                const zip = this.zip;
                if (zip) {
                    source.on("data", (chunk => {
                        if (!zip.write(chunk)) source.pause();
                    }));
                } else {
                    source.on("data", (chunk => {
                        if (!super.write(chunk)) source.pause();
                    }));
                }
            }
            pause() {
                if (this.zip) this.zip.pause();
                return super.pause();
            }
        });
        class PackSync extends Pack {
            constructor(opt) {
                super(opt);
                this[WRITEENTRYCLASS] = WriteEntrySync;
            }
            pause() {}
            resume() {}
            [STAT](job) {
                const stat = this.follow ? "statSync" : "lstatSync";
                this[ONSTAT](job, fs[stat](job.absolute));
            }
            [READDIR](job, stat) {
                this[ONREADDIR](job, fs.readdirSync(job.absolute));
            }
            [PIPE](job) {
                const source = job.entry;
                const zip = this.zip;
                if (job.readdir) {
                    job.readdir.forEach((entry => {
                        const p = this.prefix ? job.path.slice(this.prefix.length + 1) || "./" : job.path;
                        const base = p === "./" ? "" : p.replace(/\/*$/, "/");
                        this[ADDFSENTRY](base + entry);
                    }));
                }
                if (zip) {
                    source.on("data", (chunk => {
                        zip.write(chunk);
                    }));
                } else {
                    source.on("data", (chunk => {
                        super[WRITE](chunk);
                    }));
                }
            }
        }
        Pack.Sync = PackSync;
        module.exports = Pack;
    },
    6234: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const warner = __webpack_require__(8783);
        const Header = __webpack_require__(5017);
        const EE = __webpack_require__(8614);
        const Yallist = __webpack_require__(1455);
        const maxMetaEntrySize = 1024 * 1024;
        const Entry = __webpack_require__(7847);
        const Pax = __webpack_require__(9154);
        const zlib = __webpack_require__(3704);
        const gzipHeader = Buffer.from([ 31, 139 ]);
        const STATE = Symbol("state");
        const WRITEENTRY = Symbol("writeEntry");
        const READENTRY = Symbol("readEntry");
        const NEXTENTRY = Symbol("nextEntry");
        const PROCESSENTRY = Symbol("processEntry");
        const EX = Symbol("extendedHeader");
        const GEX = Symbol("globalExtendedHeader");
        const META = Symbol("meta");
        const EMITMETA = Symbol("emitMeta");
        const BUFFER = Symbol("buffer");
        const QUEUE = Symbol("queue");
        const ENDED = Symbol("ended");
        const EMITTEDEND = Symbol("emittedEnd");
        const EMIT = Symbol("emit");
        const UNZIP = Symbol("unzip");
        const CONSUMECHUNK = Symbol("consumeChunk");
        const CONSUMECHUNKSUB = Symbol("consumeChunkSub");
        const CONSUMEBODY = Symbol("consumeBody");
        const CONSUMEMETA = Symbol("consumeMeta");
        const CONSUMEHEADER = Symbol("consumeHeader");
        const CONSUMING = Symbol("consuming");
        const BUFFERCONCAT = Symbol("bufferConcat");
        const MAYBEEND = Symbol("maybeEnd");
        const WRITING = Symbol("writing");
        const ABORTED = Symbol("aborted");
        const DONE = Symbol("onDone");
        const SAW_VALID_ENTRY = Symbol("sawValidEntry");
        const SAW_NULL_BLOCK = Symbol("sawNullBlock");
        const SAW_EOF = Symbol("sawEOF");
        const noop = _ => true;
        module.exports = warner(class Parser extends EE {
            constructor(opt) {
                opt = opt || {};
                super(opt);
                this.file = opt.file || "";
                this[SAW_VALID_ENTRY] = null;
                this.on(DONE, (_ => {
                    if (this[STATE] === "begin" || this[SAW_VALID_ENTRY] === false) {
                        this.warn("TAR_BAD_ARCHIVE", "Unrecognized archive format");
                    }
                }));
                if (opt.ondone) this.on(DONE, opt.ondone); else {
                    this.on(DONE, (_ => {
                        this.emit("prefinish");
                        this.emit("finish");
                        this.emit("end");
                        this.emit("close");
                    }));
                }
                this.strict = !!opt.strict;
                this.maxMetaEntrySize = opt.maxMetaEntrySize || maxMetaEntrySize;
                this.filter = typeof opt.filter === "function" ? opt.filter : noop;
                this.writable = true;
                this.readable = false;
                this[QUEUE] = new Yallist;
                this[BUFFER] = null;
                this[READENTRY] = null;
                this[WRITEENTRY] = null;
                this[STATE] = "begin";
                this[META] = "";
                this[EX] = null;
                this[GEX] = null;
                this[ENDED] = false;
                this[UNZIP] = null;
                this[ABORTED] = false;
                this[SAW_NULL_BLOCK] = false;
                this[SAW_EOF] = false;
                if (typeof opt.onwarn === "function") this.on("warn", opt.onwarn);
                if (typeof opt.onentry === "function") this.on("entry", opt.onentry);
            }
            [CONSUMEHEADER](chunk, position) {
                if (this[SAW_VALID_ENTRY] === null) this[SAW_VALID_ENTRY] = false;
                let header;
                try {
                    header = new Header(chunk, position, this[EX], this[GEX]);
                } catch (er) {
                    return this.warn("TAR_ENTRY_INVALID", er);
                }
                if (header.nullBlock) {
                    if (this[SAW_NULL_BLOCK]) {
                        this[SAW_EOF] = true;
                        if (this[STATE] === "begin") this[STATE] = "header";
                        this[EMIT]("eof");
                    } else {
                        this[SAW_NULL_BLOCK] = true;
                        this[EMIT]("nullBlock");
                    }
                } else {
                    this[SAW_NULL_BLOCK] = false;
                    if (!header.cksumValid) this.warn("TAR_ENTRY_INVALID", "checksum failure", {
                        header
                    }); else if (!header.path) this.warn("TAR_ENTRY_INVALID", "path is required", {
                        header
                    }); else {
                        const type = header.type;
                        if (/^(Symbolic)?Link$/.test(type) && !header.linkpath) this.warn("TAR_ENTRY_INVALID", "linkpath required", {
                            header
                        }); else if (!/^(Symbolic)?Link$/.test(type) && header.linkpath) this.warn("TAR_ENTRY_INVALID", "linkpath forbidden", {
                            header
                        }); else {
                            const entry = this[WRITEENTRY] = new Entry(header, this[EX], this[GEX]);
                            if (!this[SAW_VALID_ENTRY]) {
                                if (entry.remain) {
                                    const onend = () => {
                                        if (!entry.invalid) this[SAW_VALID_ENTRY] = true;
                                    };
                                    entry.on("end", onend);
                                } else this[SAW_VALID_ENTRY] = true;
                            }
                            if (entry.meta) {
                                if (entry.size > this.maxMetaEntrySize) {
                                    entry.ignore = true;
                                    this[EMIT]("ignoredEntry", entry);
                                    this[STATE] = "ignore";
                                    entry.resume();
                                } else if (entry.size > 0) {
                                    this[META] = "";
                                    entry.on("data", (c => this[META] += c));
                                    this[STATE] = "meta";
                                }
                            } else {
                                this[EX] = null;
                                entry.ignore = entry.ignore || !this.filter(entry.path, entry);
                                if (entry.ignore) {
                                    this[EMIT]("ignoredEntry", entry);
                                    this[STATE] = entry.remain ? "ignore" : "header";
                                    entry.resume();
                                } else {
                                    if (entry.remain) this[STATE] = "body"; else {
                                        this[STATE] = "header";
                                        entry.end();
                                    }
                                    if (!this[READENTRY]) {
                                        this[QUEUE].push(entry);
                                        this[NEXTENTRY]();
                                    } else this[QUEUE].push(entry);
                                }
                            }
                        }
                    }
                }
            }
            [PROCESSENTRY](entry) {
                let go = true;
                if (!entry) {
                    this[READENTRY] = null;
                    go = false;
                } else if (Array.isArray(entry)) this.emit.apply(this, entry); else {
                    this[READENTRY] = entry;
                    this.emit("entry", entry);
                    if (!entry.emittedEnd) {
                        entry.on("end", (_ => this[NEXTENTRY]()));
                        go = false;
                    }
                }
                return go;
            }
            [NEXTENTRY]() {
                do {} while (this[PROCESSENTRY](this[QUEUE].shift()));
                if (!this[QUEUE].length) {
                    const re = this[READENTRY];
                    const drainNow = !re || re.flowing || re.size === re.remain;
                    if (drainNow) {
                        if (!this[WRITING]) this.emit("drain");
                    } else re.once("drain", (_ => this.emit("drain")));
                }
            }
            [CONSUMEBODY](chunk, position) {
                const entry = this[WRITEENTRY];
                const br = entry.blockRemain;
                const c = br >= chunk.length && position === 0 ? chunk : chunk.slice(position, position + br);
                entry.write(c);
                if (!entry.blockRemain) {
                    this[STATE] = "header";
                    this[WRITEENTRY] = null;
                    entry.end();
                }
                return c.length;
            }
            [CONSUMEMETA](chunk, position) {
                const entry = this[WRITEENTRY];
                const ret = this[CONSUMEBODY](chunk, position);
                if (!this[WRITEENTRY]) this[EMITMETA](entry);
                return ret;
            }
            [EMIT](ev, data, extra) {
                if (!this[QUEUE].length && !this[READENTRY]) this.emit(ev, data, extra); else this[QUEUE].push([ ev, data, extra ]);
            }
            [EMITMETA](entry) {
                this[EMIT]("meta", this[META]);
                switch (entry.type) {
                  case "ExtendedHeader":
                  case "OldExtendedHeader":
                    this[EX] = Pax.parse(this[META], this[EX], false);
                    break;

                  case "GlobalExtendedHeader":
                    this[GEX] = Pax.parse(this[META], this[GEX], true);
                    break;

                  case "NextFileHasLongPath":
                  case "OldGnuLongPath":
                    this[EX] = this[EX] || Object.create(null);
                    this[EX].path = this[META].replace(/\0.*/, "");
                    break;

                  case "NextFileHasLongLinkpath":
                    this[EX] = this[EX] || Object.create(null);
                    this[EX].linkpath = this[META].replace(/\0.*/, "");
                    break;

                  default:
                    throw new Error("unknown meta: " + entry.type);
                }
            }
            abort(error) {
                this[ABORTED] = true;
                this.emit("abort", error);
                this.warn("TAR_ABORT", error, {
                    recoverable: false
                });
            }
            write(chunk) {
                if (this[ABORTED]) return;
                if (this[UNZIP] === null && chunk) {
                    if (this[BUFFER]) {
                        chunk = Buffer.concat([ this[BUFFER], chunk ]);
                        this[BUFFER] = null;
                    }
                    if (chunk.length < gzipHeader.length) {
                        this[BUFFER] = chunk;
                        return true;
                    }
                    for (let i = 0; this[UNZIP] === null && i < gzipHeader.length; i++) {
                        if (chunk[i] !== gzipHeader[i]) this[UNZIP] = false;
                    }
                    if (this[UNZIP] === null) {
                        const ended = this[ENDED];
                        this[ENDED] = false;
                        this[UNZIP] = new zlib.Unzip;
                        this[UNZIP].on("data", (chunk => this[CONSUMECHUNK](chunk)));
                        this[UNZIP].on("error", (er => this.abort(er)));
                        this[UNZIP].on("end", (_ => {
                            this[ENDED] = true;
                            this[CONSUMECHUNK]();
                        }));
                        this[WRITING] = true;
                        const ret = this[UNZIP][ended ? "end" : "write"](chunk);
                        this[WRITING] = false;
                        return ret;
                    }
                }
                this[WRITING] = true;
                if (this[UNZIP]) this[UNZIP].write(chunk); else this[CONSUMECHUNK](chunk);
                this[WRITING] = false;
                const ret = this[QUEUE].length ? false : this[READENTRY] ? this[READENTRY].flowing : true;
                if (!ret && !this[QUEUE].length) this[READENTRY].once("drain", (_ => this.emit("drain")));
                return ret;
            }
            [BUFFERCONCAT](c) {
                if (c && !this[ABORTED]) this[BUFFER] = this[BUFFER] ? Buffer.concat([ this[BUFFER], c ]) : c;
            }
            [MAYBEEND]() {
                if (this[ENDED] && !this[EMITTEDEND] && !this[ABORTED] && !this[CONSUMING]) {
                    this[EMITTEDEND] = true;
                    const entry = this[WRITEENTRY];
                    if (entry && entry.blockRemain) {
                        const have = this[BUFFER] ? this[BUFFER].length : 0;
                        this.warn("TAR_BAD_ARCHIVE", `Truncated input (needed ${entry.blockRemain} more bytes, only ${have} available)`, {
                            entry
                        });
                        if (this[BUFFER]) entry.write(this[BUFFER]);
                        entry.end();
                    }
                    this[EMIT](DONE);
                }
            }
            [CONSUMECHUNK](chunk) {
                if (this[CONSUMING]) this[BUFFERCONCAT](chunk); else if (!chunk && !this[BUFFER]) this[MAYBEEND](); else {
                    this[CONSUMING] = true;
                    if (this[BUFFER]) {
                        this[BUFFERCONCAT](chunk);
                        const c = this[BUFFER];
                        this[BUFFER] = null;
                        this[CONSUMECHUNKSUB](c);
                    } else this[CONSUMECHUNKSUB](chunk);
                    while (this[BUFFER] && this[BUFFER].length >= 512 && !this[ABORTED] && !this[SAW_EOF]) {
                        const c = this[BUFFER];
                        this[BUFFER] = null;
                        this[CONSUMECHUNKSUB](c);
                    }
                    this[CONSUMING] = false;
                }
                if (!this[BUFFER] || this[ENDED]) this[MAYBEEND]();
            }
            [CONSUMECHUNKSUB](chunk) {
                let position = 0;
                const length = chunk.length;
                while (position + 512 <= length && !this[ABORTED] && !this[SAW_EOF]) {
                    switch (this[STATE]) {
                      case "begin":
                      case "header":
                        this[CONSUMEHEADER](chunk, position);
                        position += 512;
                        break;

                      case "ignore":
                      case "body":
                        position += this[CONSUMEBODY](chunk, position);
                        break;

                      case "meta":
                        position += this[CONSUMEMETA](chunk, position);
                        break;

                      default:
                        throw new Error("invalid state: " + this[STATE]);
                    }
                }
                if (position < length) {
                    if (this[BUFFER]) this[BUFFER] = Buffer.concat([ chunk.slice(position), this[BUFFER] ]); else this[BUFFER] = chunk.slice(position);
                }
            }
            end(chunk) {
                if (!this[ABORTED]) {
                    if (this[UNZIP]) this[UNZIP].end(chunk); else {
                        this[ENDED] = true;
                        this.write(chunk);
                    }
                }
            }
        });
    },
    7119: (module, __unused_webpack_exports, __webpack_require__) => {
        const assert = __webpack_require__(2357);
        module.exports = () => {
            const queues = new Map;
            const reservations = new Map;
            const {join} = __webpack_require__(5622);
            const getDirs = path => join(path).split(/[\\/]/).slice(0, -1).reduce(((set, path) => set.length ? set.concat(join(set[set.length - 1], path)) : [ path ]), []);
            const running = new Set;
            const getQueues = fn => {
                const res = reservations.get(fn);
                if (!res) throw new Error("function does not have any path reservations");
                return {
                    paths: res.paths.map((path => queues.get(path))),
                    dirs: [ ...res.dirs ].map((path => queues.get(path)))
                };
            };
            const check = fn => {
                const {paths, dirs} = getQueues(fn);
                return paths.every((q => q[0] === fn)) && dirs.every((q => q[0] instanceof Set && q[0].has(fn)));
            };
            const run = fn => {
                if (running.has(fn) || !check(fn)) return false;
                running.add(fn);
                fn((() => clear(fn)));
                return true;
            };
            const clear = fn => {
                if (!running.has(fn)) return false;
                const {paths, dirs} = reservations.get(fn);
                const next = new Set;
                paths.forEach((path => {
                    const q = queues.get(path);
                    assert.equal(q[0], fn);
                    if (q.length === 1) queues.delete(path); else {
                        q.shift();
                        if (typeof q[0] === "function") next.add(q[0]); else q[0].forEach((fn => next.add(fn)));
                    }
                }));
                dirs.forEach((dir => {
                    const q = queues.get(dir);
                    assert(q[0] instanceof Set);
                    if (q[0].size === 1 && q.length === 1) queues.delete(dir); else if (q[0].size === 1) {
                        q.shift();
                        next.add(q[0]);
                    } else q[0].delete(fn);
                }));
                running.delete(fn);
                next.forEach((fn => run(fn)));
                return true;
            };
            const reserve = (paths, fn) => {
                const dirs = new Set(paths.map((path => getDirs(path))).reduce(((a, b) => a.concat(b))));
                reservations.set(fn, {
                    dirs,
                    paths
                });
                paths.forEach((path => {
                    const q = queues.get(path);
                    if (!q) queues.set(path, [ fn ]); else q.push(fn);
                }));
                dirs.forEach((dir => {
                    const q = queues.get(dir);
                    if (!q) queues.set(dir, [ new Set([ fn ]) ]); else if (q[q.length - 1] instanceof Set) q[q.length - 1].add(fn); else q.push(new Set([ fn ]));
                }));
                return run(fn);
            };
            return {
                check,
                reserve
            };
        };
    },
    9154: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const Header = __webpack_require__(5017);
        const path = __webpack_require__(5622);
        class Pax {
            constructor(obj, global) {
                this.atime = obj.atime || null;
                this.charset = obj.charset || null;
                this.comment = obj.comment || null;
                this.ctime = obj.ctime || null;
                this.gid = obj.gid || null;
                this.gname = obj.gname || null;
                this.linkpath = obj.linkpath || null;
                this.mtime = obj.mtime || null;
                this.path = obj.path || null;
                this.size = obj.size || null;
                this.uid = obj.uid || null;
                this.uname = obj.uname || null;
                this.dev = obj.dev || null;
                this.ino = obj.ino || null;
                this.nlink = obj.nlink || null;
                this.global = global || false;
            }
            encode() {
                const body = this.encodeBody();
                if (body === "") return null;
                const bodyLen = Buffer.byteLength(body);
                const bufLen = 512 * Math.ceil(1 + bodyLen / 512);
                const buf = Buffer.allocUnsafe(bufLen);
                for (let i = 0; i < 512; i++) buf[i] = 0;
                new Header({
                    path: ("PaxHeader/" + path.basename(this.path)).slice(0, 99),
                    mode: this.mode || 420,
                    uid: this.uid || null,
                    gid: this.gid || null,
                    size: bodyLen,
                    mtime: this.mtime || null,
                    type: this.global ? "GlobalExtendedHeader" : "ExtendedHeader",
                    linkpath: "",
                    uname: this.uname || "",
                    gname: this.gname || "",
                    devmaj: 0,
                    devmin: 0,
                    atime: this.atime || null,
                    ctime: this.ctime || null
                }).encode(buf);
                buf.write(body, 512, bodyLen, "utf8");
                for (let i = bodyLen + 512; i < buf.length; i++) buf[i] = 0;
                return buf;
            }
            encodeBody() {
                return this.encodeField("path") + this.encodeField("ctime") + this.encodeField("atime") + this.encodeField("dev") + this.encodeField("ino") + this.encodeField("nlink") + this.encodeField("charset") + this.encodeField("comment") + this.encodeField("gid") + this.encodeField("gname") + this.encodeField("linkpath") + this.encodeField("mtime") + this.encodeField("size") + this.encodeField("uid") + this.encodeField("uname");
            }
            encodeField(field) {
                if (this[field] === null || this[field] === undefined) return "";
                const v = this[field] instanceof Date ? this[field].getTime() / 1e3 : this[field];
                const s = " " + (field === "dev" || field === "ino" || field === "nlink" ? "SCHILY." : "") + field + "=" + v + "\n";
                const byteLen = Buffer.byteLength(s);
                let digits = Math.floor(Math.log(byteLen) / Math.log(10)) + 1;
                if (byteLen + digits >= Math.pow(10, digits)) digits += 1;
                const len = digits + byteLen;
                return len + s;
            }
        }
        Pax.parse = (string, ex, g) => new Pax(merge(parseKV(string), ex), g);
        const merge = (a, b) => b ? Object.keys(a).reduce(((s, k) => (s[k] = a[k], s)), b) : a;
        const parseKV = string => string.replace(/\n$/, "").split("\n").reduce(parseKVLine, Object.create(null));
        const parseKVLine = (set, line) => {
            const n = parseInt(line, 10);
            if (n !== Buffer.byteLength(line) + 1) return set;
            line = line.substr((n + " ").length);
            const kv = line.split("=");
            const k = kv.shift().replace(/^SCHILY\.(dev|ino|nlink)/, "$1");
            if (!k) return set;
            const v = kv.join("=");
            set[k] = /^([A-Z]+\.)?([mac]|birth|creation)time$/.test(k) ? new Date(v * 1e3) : /^[0-9]+$/.test(v) ? +v : v;
            return set;
        };
        module.exports = Pax;
    },
    7847: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const MiniPass = __webpack_require__(2253);
        const SLURP = Symbol("slurp");
        module.exports = class ReadEntry extends MiniPass {
            constructor(header, ex, gex) {
                super();
                this.pause();
                this.extended = ex;
                this.globalExtended = gex;
                this.header = header;
                this.startBlockSize = 512 * Math.ceil(header.size / 512);
                this.blockRemain = this.startBlockSize;
                this.remain = header.size;
                this.type = header.type;
                this.meta = false;
                this.ignore = false;
                switch (this.type) {
                  case "File":
                  case "OldFile":
                  case "Link":
                  case "SymbolicLink":
                  case "CharacterDevice":
                  case "BlockDevice":
                  case "Directory":
                  case "FIFO":
                  case "ContiguousFile":
                  case "GNUDumpDir":
                    break;

                  case "NextFileHasLongLinkpath":
                  case "NextFileHasLongPath":
                  case "OldGnuLongPath":
                  case "GlobalExtendedHeader":
                  case "ExtendedHeader":
                  case "OldExtendedHeader":
                    this.meta = true;
                    break;

                  default:
                    this.ignore = true;
                }
                this.path = header.path;
                this.mode = header.mode;
                if (this.mode) this.mode = this.mode & 4095;
                this.uid = header.uid;
                this.gid = header.gid;
                this.uname = header.uname;
                this.gname = header.gname;
                this.size = header.size;
                this.mtime = header.mtime;
                this.atime = header.atime;
                this.ctime = header.ctime;
                this.linkpath = header.linkpath;
                this.uname = header.uname;
                this.gname = header.gname;
                if (ex) this[SLURP](ex);
                if (gex) this[SLURP](gex, true);
            }
            write(data) {
                const writeLen = data.length;
                if (writeLen > this.blockRemain) throw new Error("writing more to entry than is appropriate");
                const r = this.remain;
                const br = this.blockRemain;
                this.remain = Math.max(0, r - writeLen);
                this.blockRemain = Math.max(0, br - writeLen);
                if (this.ignore) return true;
                if (r >= writeLen) return super.write(data);
                return super.write(data.slice(0, r));
            }
            [SLURP](ex, global) {
                for (const k in ex) {
                    if (ex[k] !== null && ex[k] !== undefined && !(global && k === "path")) this[k] = ex[k];
                }
            }
        };
    },
    3666: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const hlo = __webpack_require__(7461);
        const Pack = __webpack_require__(5843);
        const fs = __webpack_require__(5747);
        const fsm = __webpack_require__(8553);
        const t = __webpack_require__(1090);
        const path = __webpack_require__(5622);
        const Header = __webpack_require__(5017);
        module.exports = (opt_, files, cb) => {
            const opt = hlo(opt_);
            if (!opt.file) throw new TypeError("file is required");
            if (opt.gzip) throw new TypeError("cannot append to compressed archives");
            if (!files || !Array.isArray(files) || !files.length) throw new TypeError("no files or directories specified");
            files = Array.from(files);
            return opt.sync ? replaceSync(opt, files) : replace(opt, files, cb);
        };
        const replaceSync = (opt, files) => {
            const p = new Pack.Sync(opt);
            let threw = true;
            let fd;
            let position;
            try {
                try {
                    fd = fs.openSync(opt.file, "r+");
                } catch (er) {
                    if (er.code === "ENOENT") fd = fs.openSync(opt.file, "w+"); else throw er;
                }
                const st = fs.fstatSync(fd);
                const headBuf = Buffer.alloc(512);
                POSITION: for (position = 0; position < st.size; position += 512) {
                    for (let bufPos = 0, bytes = 0; bufPos < 512; bufPos += bytes) {
                        bytes = fs.readSync(fd, headBuf, bufPos, headBuf.length - bufPos, position + bufPos);
                        if (position === 0 && headBuf[0] === 31 && headBuf[1] === 139) throw new Error("cannot append to compressed archives");
                        if (!bytes) break POSITION;
                    }
                    const h = new Header(headBuf);
                    if (!h.cksumValid) break;
                    const entryBlockSize = 512 * Math.ceil(h.size / 512);
                    if (position + entryBlockSize + 512 > st.size) break;
                    position += entryBlockSize;
                    if (opt.mtimeCache) opt.mtimeCache.set(h.path, h.mtime);
                }
                threw = false;
                streamSync(opt, p, position, fd, files);
            } finally {
                if (threw) {
                    try {
                        fs.closeSync(fd);
                    } catch (er) {}
                }
            }
        };
        const streamSync = (opt, p, position, fd, files) => {
            const stream = new fsm.WriteStreamSync(opt.file, {
                fd,
                start: position
            });
            p.pipe(stream);
            addFilesSync(p, files);
        };
        const replace = (opt, files, cb) => {
            files = Array.from(files);
            const p = new Pack(opt);
            const getPos = (fd, size, cb_) => {
                const cb = (er, pos) => {
                    if (er) fs.close(fd, (_ => cb_(er))); else cb_(null, pos);
                };
                let position = 0;
                if (size === 0) return cb(null, 0);
                let bufPos = 0;
                const headBuf = Buffer.alloc(512);
                const onread = (er, bytes) => {
                    if (er) return cb(er);
                    bufPos += bytes;
                    if (bufPos < 512 && bytes) {
                        return fs.read(fd, headBuf, bufPos, headBuf.length - bufPos, position + bufPos, onread);
                    }
                    if (position === 0 && headBuf[0] === 31 && headBuf[1] === 139) return cb(new Error("cannot append to compressed archives"));
                    if (bufPos < 512) return cb(null, position);
                    const h = new Header(headBuf);
                    if (!h.cksumValid) return cb(null, position);
                    const entryBlockSize = 512 * Math.ceil(h.size / 512);
                    if (position + entryBlockSize + 512 > size) return cb(null, position);
                    position += entryBlockSize + 512;
                    if (position >= size) return cb(null, position);
                    if (opt.mtimeCache) opt.mtimeCache.set(h.path, h.mtime);
                    bufPos = 0;
                    fs.read(fd, headBuf, 0, 512, position, onread);
                };
                fs.read(fd, headBuf, 0, 512, position, onread);
            };
            const promise = new Promise(((resolve, reject) => {
                p.on("error", reject);
                let flag = "r+";
                const onopen = (er, fd) => {
                    if (er && er.code === "ENOENT" && flag === "r+") {
                        flag = "w+";
                        return fs.open(opt.file, flag, onopen);
                    }
                    if (er) return reject(er);
                    fs.fstat(fd, ((er, st) => {
                        if (er) return reject(er);
                        getPos(fd, st.size, ((er, position) => {
                            if (er) return reject(er);
                            const stream = new fsm.WriteStream(opt.file, {
                                fd,
                                start: position
                            });
                            p.pipe(stream);
                            stream.on("error", reject);
                            stream.on("close", resolve);
                            addFilesAsync(p, files);
                        }));
                    }));
                };
                fs.open(opt.file, flag, onopen);
            }));
            return cb ? promise.then(cb, cb) : promise;
        };
        const addFilesSync = (p, files) => {
            files.forEach((file => {
                if (file.charAt(0) === "@") {
                    t({
                        file: path.resolve(p.cwd, file.substr(1)),
                        sync: true,
                        noResume: true,
                        onentry: entry => p.add(entry)
                    });
                } else p.add(file);
            }));
            p.end();
        };
        const addFilesAsync = (p, files) => {
            while (files.length) {
                const file = files.shift();
                if (file.charAt(0) === "@") {
                    return t({
                        file: path.resolve(p.cwd, file.substr(1)),
                        noResume: true,
                        onentry: entry => p.add(entry)
                    }).then((_ => addFilesAsync(p, files)));
                } else p.add(file);
            }
            p.end();
        };
    },
    9806: (__unused_webpack_module, exports) => {
        "use strict";
        exports.name = new Map([ [ "0", "File" ], [ "", "OldFile" ], [ "1", "Link" ], [ "2", "SymbolicLink" ], [ "3", "CharacterDevice" ], [ "4", "BlockDevice" ], [ "5", "Directory" ], [ "6", "FIFO" ], [ "7", "ContiguousFile" ], [ "g", "GlobalExtendedHeader" ], [ "x", "ExtendedHeader" ], [ "A", "SolarisACL" ], [ "D", "GNUDumpDir" ], [ "I", "Inode" ], [ "K", "NextFileHasLongLinkpath" ], [ "L", "NextFileHasLongPath" ], [ "M", "ContinuationFile" ], [ "N", "OldGnuLongPath" ], [ "S", "SparseFile" ], [ "V", "TapeVolumeHeader" ], [ "X", "OldExtendedHeader" ] ]);
        exports.code = new Map(Array.from(exports.name).map((kv => [ kv[1], kv[0] ])));
    },
    2864: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const assert = __webpack_require__(2357);
        const Parser = __webpack_require__(6234);
        const fs = __webpack_require__(5747);
        const fsm = __webpack_require__(8553);
        const path = __webpack_require__(5622);
        const mkdir = __webpack_require__(3956);
        const wc = __webpack_require__(6564);
        const pathReservations = __webpack_require__(7119);
        const ONENTRY = Symbol("onEntry");
        const CHECKFS = Symbol("checkFs");
        const CHECKFS2 = Symbol("checkFs2");
        const ISREUSABLE = Symbol("isReusable");
        const MAKEFS = Symbol("makeFs");
        const FILE = Symbol("file");
        const DIRECTORY = Symbol("directory");
        const LINK = Symbol("link");
        const SYMLINK = Symbol("symlink");
        const HARDLINK = Symbol("hardlink");
        const UNSUPPORTED = Symbol("unsupported");
        const CHECKPATH = Symbol("checkPath");
        const MKDIR = Symbol("mkdir");
        const ONERROR = Symbol("onError");
        const PENDING = Symbol("pending");
        const PEND = Symbol("pend");
        const UNPEND = Symbol("unpend");
        const ENDED = Symbol("ended");
        const MAYBECLOSE = Symbol("maybeClose");
        const SKIP = Symbol("skip");
        const DOCHOWN = Symbol("doChown");
        const UID = Symbol("uid");
        const GID = Symbol("gid");
        const crypto = __webpack_require__(6417);
        const getFlag = __webpack_require__(8512);
        const neverCalled = () => {
            throw new Error("sync function called cb somehow?!?");
        };
        const unlinkFile = (path, cb) => {
            if (process.platform !== "win32") return fs.unlink(path, cb);
            const name = path + ".DELETE." + crypto.randomBytes(16).toString("hex");
            fs.rename(path, name, (er => {
                if (er) return cb(er);
                fs.unlink(name, cb);
            }));
        };
        const unlinkFileSync = path => {
            if (process.platform !== "win32") return fs.unlinkSync(path);
            const name = path + ".DELETE." + crypto.randomBytes(16).toString("hex");
            fs.renameSync(path, name);
            fs.unlinkSync(name);
        };
        const uint32 = (a, b, c) => a === a >>> 0 ? a : b === b >>> 0 ? b : c;
        class Unpack extends Parser {
            constructor(opt) {
                if (!opt) opt = {};
                opt.ondone = _ => {
                    this[ENDED] = true;
                    this[MAYBECLOSE]();
                };
                super(opt);
                this.reservations = pathReservations();
                this.transform = typeof opt.transform === "function" ? opt.transform : null;
                this.writable = true;
                this.readable = false;
                this[PENDING] = 0;
                this[ENDED] = false;
                this.dirCache = opt.dirCache || new Map;
                if (typeof opt.uid === "number" || typeof opt.gid === "number") {
                    if (typeof opt.uid !== "number" || typeof opt.gid !== "number") throw new TypeError("cannot set owner without number uid and gid");
                    if (opt.preserveOwner) {
                        throw new TypeError("cannot preserve owner in archive and also set owner explicitly");
                    }
                    this.uid = opt.uid;
                    this.gid = opt.gid;
                    this.setOwner = true;
                } else {
                    this.uid = null;
                    this.gid = null;
                    this.setOwner = false;
                }
                if (opt.preserveOwner === undefined && typeof opt.uid !== "number") this.preserveOwner = process.getuid && process.getuid() === 0; else this.preserveOwner = !!opt.preserveOwner;
                this.processUid = (this.preserveOwner || this.setOwner) && process.getuid ? process.getuid() : null;
                this.processGid = (this.preserveOwner || this.setOwner) && process.getgid ? process.getgid() : null;
                this.forceChown = opt.forceChown === true;
                this.win32 = !!opt.win32 || process.platform === "win32";
                this.newer = !!opt.newer;
                this.keep = !!opt.keep;
                this.noMtime = !!opt.noMtime;
                this.preservePaths = !!opt.preservePaths;
                this.unlink = !!opt.unlink;
                this.cwd = path.resolve(opt.cwd || process.cwd());
                this.strip = +opt.strip || 0;
                this.processUmask = opt.noChmod ? 0 : process.umask();
                this.umask = typeof opt.umask === "number" ? opt.umask : this.processUmask;
                this.dmode = opt.dmode || 511 & ~this.umask;
                this.fmode = opt.fmode || 438 & ~this.umask;
                this.on("entry", (entry => this[ONENTRY](entry)));
            }
            warn(code, msg, data = {}) {
                if (code === "TAR_BAD_ARCHIVE" || code === "TAR_ABORT") data.recoverable = false;
                return super.warn(code, msg, data);
            }
            [MAYBECLOSE]() {
                if (this[ENDED] && this[PENDING] === 0) {
                    this.emit("prefinish");
                    this.emit("finish");
                    this.emit("end");
                    this.emit("close");
                }
            }
            [CHECKPATH](entry) {
                if (this.strip) {
                    const parts = entry.path.split(/\/|\\/);
                    if (parts.length < this.strip) return false;
                    entry.path = parts.slice(this.strip).join("/");
                    if (entry.type === "Link") {
                        const linkparts = entry.linkpath.split(/\/|\\/);
                        if (linkparts.length >= this.strip) entry.linkpath = linkparts.slice(this.strip).join("/");
                    }
                }
                if (!this.preservePaths) {
                    const p = entry.path;
                    if (p.match(/(^|\/|\\)\.\.(\\|\/|$)/)) {
                        this.warn("TAR_ENTRY_ERROR", `path contains '..'`, {
                            entry,
                            path: p
                        });
                        return false;
                    }
                    if (path.win32.isAbsolute(p)) {
                        const parsed = path.win32.parse(p);
                        entry.path = p.substr(parsed.root.length);
                        const r = parsed.root;
                        this.warn("TAR_ENTRY_INFO", `stripping ${r} from absolute path`, {
                            entry,
                            path: p
                        });
                    }
                }
                if (this.win32) {
                    const parsed = path.win32.parse(entry.path);
                    entry.path = parsed.root === "" ? wc.encode(entry.path) : parsed.root + wc.encode(entry.path.substr(parsed.root.length));
                }
                if (path.isAbsolute(entry.path)) entry.absolute = entry.path; else entry.absolute = path.resolve(this.cwd, entry.path);
                return true;
            }
            [ONENTRY](entry) {
                if (!this[CHECKPATH](entry)) return entry.resume();
                assert.equal(typeof entry.absolute, "string");
                switch (entry.type) {
                  case "Directory":
                  case "GNUDumpDir":
                    if (entry.mode) entry.mode = entry.mode | 448;

                  case "File":
                  case "OldFile":
                  case "ContiguousFile":
                  case "Link":
                  case "SymbolicLink":
                    return this[CHECKFS](entry);

                  case "CharacterDevice":
                  case "BlockDevice":
                  case "FIFO":
                  default:
                    return this[UNSUPPORTED](entry);
                }
            }
            [ONERROR](er, entry) {
                if (er.name === "CwdError") this.emit("error", er); else {
                    this.warn("TAR_ENTRY_ERROR", er, {
                        entry
                    });
                    this[UNPEND]();
                    entry.resume();
                }
            }
            [MKDIR](dir, mode, cb) {
                mkdir(dir, {
                    uid: this.uid,
                    gid: this.gid,
                    processUid: this.processUid,
                    processGid: this.processGid,
                    umask: this.processUmask,
                    preserve: this.preservePaths,
                    unlink: this.unlink,
                    cache: this.dirCache,
                    cwd: this.cwd,
                    mode,
                    noChmod: this.noChmod
                }, cb);
            }
            [DOCHOWN](entry) {
                return this.forceChown || this.preserveOwner && (typeof entry.uid === "number" && entry.uid !== this.processUid || typeof entry.gid === "number" && entry.gid !== this.processGid) || (typeof this.uid === "number" && this.uid !== this.processUid || typeof this.gid === "number" && this.gid !== this.processGid);
            }
            [UID](entry) {
                return uint32(this.uid, entry.uid, this.processUid);
            }
            [GID](entry) {
                return uint32(this.gid, entry.gid, this.processGid);
            }
            [FILE](entry, fullyDone) {
                const mode = entry.mode & 4095 || this.fmode;
                const stream = new fsm.WriteStream(entry.absolute, {
                    flags: getFlag(entry.size),
                    mode,
                    autoClose: false
                });
                stream.on("error", (er => this[ONERROR](er, entry)));
                let actions = 1;
                const done = er => {
                    if (er) return this[ONERROR](er, entry);
                    if (--actions === 0) {
                        fs.close(stream.fd, (er => {
                            fullyDone();
                            er ? this[ONERROR](er, entry) : this[UNPEND]();
                        }));
                    }
                };
                stream.on("finish", (_ => {
                    const abs = entry.absolute;
                    const fd = stream.fd;
                    if (entry.mtime && !this.noMtime) {
                        actions++;
                        const atime = entry.atime || new Date;
                        const mtime = entry.mtime;
                        fs.futimes(fd, atime, mtime, (er => er ? fs.utimes(abs, atime, mtime, (er2 => done(er2 && er))) : done()));
                    }
                    if (this[DOCHOWN](entry)) {
                        actions++;
                        const uid = this[UID](entry);
                        const gid = this[GID](entry);
                        fs.fchown(fd, uid, gid, (er => er ? fs.chown(abs, uid, gid, (er2 => done(er2 && er))) : done()));
                    }
                    done();
                }));
                const tx = this.transform ? this.transform(entry) || entry : entry;
                if (tx !== entry) {
                    tx.on("error", (er => this[ONERROR](er, entry)));
                    entry.pipe(tx);
                }
                tx.pipe(stream);
            }
            [DIRECTORY](entry, fullyDone) {
                const mode = entry.mode & 4095 || this.dmode;
                this[MKDIR](entry.absolute, mode, (er => {
                    if (er) {
                        fullyDone();
                        return this[ONERROR](er, entry);
                    }
                    let actions = 1;
                    const done = _ => {
                        if (--actions === 0) {
                            fullyDone();
                            this[UNPEND]();
                            entry.resume();
                        }
                    };
                    if (entry.mtime && !this.noMtime) {
                        actions++;
                        fs.utimes(entry.absolute, entry.atime || new Date, entry.mtime, done);
                    }
                    if (this[DOCHOWN](entry)) {
                        actions++;
                        fs.chown(entry.absolute, this[UID](entry), this[GID](entry), done);
                    }
                    done();
                }));
            }
            [UNSUPPORTED](entry) {
                entry.unsupported = true;
                this.warn("TAR_ENTRY_UNSUPPORTED", `unsupported entry type: ${entry.type}`, {
                    entry
                });
                entry.resume();
            }
            [SYMLINK](entry, done) {
                this[LINK](entry, entry.linkpath, "symlink", done);
            }
            [HARDLINK](entry, done) {
                this[LINK](entry, path.resolve(this.cwd, entry.linkpath), "link", done);
            }
            [PEND]() {
                this[PENDING]++;
            }
            [UNPEND]() {
                this[PENDING]--;
                this[MAYBECLOSE]();
            }
            [SKIP](entry) {
                this[UNPEND]();
                entry.resume();
            }
            [ISREUSABLE](entry, st) {
                return entry.type === "File" && !this.unlink && st.isFile() && st.nlink <= 1 && process.platform !== "win32";
            }
            [CHECKFS](entry) {
                this[PEND]();
                const paths = [ entry.path ];
                if (entry.linkpath) paths.push(entry.linkpath);
                this.reservations.reserve(paths, (done => this[CHECKFS2](entry, done)));
            }
            [CHECKFS2](entry, done) {
                this[MKDIR](path.dirname(entry.absolute), this.dmode, (er => {
                    if (er) {
                        done();
                        return this[ONERROR](er, entry);
                    }
                    fs.lstat(entry.absolute, ((er, st) => {
                        if (st && (this.keep || this.newer && st.mtime > entry.mtime)) {
                            this[SKIP](entry);
                            done();
                        } else if (er || this[ISREUSABLE](entry, st)) this[MAKEFS](null, entry, done); else if (st.isDirectory()) {
                            if (entry.type === "Directory") {
                                if (!this.noChmod && (!entry.mode || (st.mode & 4095) === entry.mode)) this[MAKEFS](null, entry, done); else {
                                    fs.chmod(entry.absolute, entry.mode, (er => this[MAKEFS](er, entry, done)));
                                }
                            } else fs.rmdir(entry.absolute, (er => this[MAKEFS](er, entry, done)));
                        } else unlinkFile(entry.absolute, (er => this[MAKEFS](er, entry, done)));
                    }));
                }));
            }
            [MAKEFS](er, entry, done) {
                if (er) return this[ONERROR](er, entry);
                switch (entry.type) {
                  case "File":
                  case "OldFile":
                  case "ContiguousFile":
                    return this[FILE](entry, done);

                  case "Link":
                    return this[HARDLINK](entry, done);

                  case "SymbolicLink":
                    return this[SYMLINK](entry, done);

                  case "Directory":
                  case "GNUDumpDir":
                    return this[DIRECTORY](entry, done);
                }
            }
            [LINK](entry, linkpath, link, done) {
                fs[link](linkpath, entry.absolute, (er => {
                    if (er) return this[ONERROR](er, entry);
                    done();
                    this[UNPEND]();
                    entry.resume();
                }));
            }
        }
        class UnpackSync extends Unpack {
            [CHECKFS](entry) {
                const er = this[MKDIR](path.dirname(entry.absolute), this.dmode, neverCalled);
                if (er) return this[ONERROR](er, entry);
                try {
                    const st = fs.lstatSync(entry.absolute);
                    if (this.keep || this.newer && st.mtime > entry.mtime) return this[SKIP](entry); else if (this[ISREUSABLE](entry, st)) return this[MAKEFS](null, entry, neverCalled); else {
                        try {
                            if (st.isDirectory()) {
                                if (entry.type === "Directory") {
                                    if (!this.noChmod && entry.mode && (st.mode & 4095) !== entry.mode) fs.chmodSync(entry.absolute, entry.mode);
                                } else fs.rmdirSync(entry.absolute);
                            } else unlinkFileSync(entry.absolute);
                            return this[MAKEFS](null, entry, neverCalled);
                        } catch (er) {
                            return this[ONERROR](er, entry);
                        }
                    }
                } catch (er) {
                    return this[MAKEFS](null, entry, neverCalled);
                }
            }
            [FILE](entry, _) {
                const mode = entry.mode & 4095 || this.fmode;
                const oner = er => {
                    let closeError;
                    try {
                        fs.closeSync(fd);
                    } catch (e) {
                        closeError = e;
                    }
                    if (er || closeError) this[ONERROR](er || closeError, entry);
                };
                let fd;
                try {
                    fd = fs.openSync(entry.absolute, getFlag(entry.size), mode);
                } catch (er) {
                    return oner(er);
                }
                const tx = this.transform ? this.transform(entry) || entry : entry;
                if (tx !== entry) {
                    tx.on("error", (er => this[ONERROR](er, entry)));
                    entry.pipe(tx);
                }
                tx.on("data", (chunk => {
                    try {
                        fs.writeSync(fd, chunk, 0, chunk.length);
                    } catch (er) {
                        oner(er);
                    }
                }));
                tx.on("end", (_ => {
                    let er = null;
                    if (entry.mtime && !this.noMtime) {
                        const atime = entry.atime || new Date;
                        const mtime = entry.mtime;
                        try {
                            fs.futimesSync(fd, atime, mtime);
                        } catch (futimeser) {
                            try {
                                fs.utimesSync(entry.absolute, atime, mtime);
                            } catch (utimeser) {
                                er = futimeser;
                            }
                        }
                    }
                    if (this[DOCHOWN](entry)) {
                        const uid = this[UID](entry);
                        const gid = this[GID](entry);
                        try {
                            fs.fchownSync(fd, uid, gid);
                        } catch (fchowner) {
                            try {
                                fs.chownSync(entry.absolute, uid, gid);
                            } catch (chowner) {
                                er = er || fchowner;
                            }
                        }
                    }
                    oner(er);
                }));
            }
            [DIRECTORY](entry, _) {
                const mode = entry.mode & 4095 || this.dmode;
                const er = this[MKDIR](entry.absolute, mode);
                if (er) return this[ONERROR](er, entry);
                if (entry.mtime && !this.noMtime) {
                    try {
                        fs.utimesSync(entry.absolute, entry.atime || new Date, entry.mtime);
                    } catch (er) {}
                }
                if (this[DOCHOWN](entry)) {
                    try {
                        fs.chownSync(entry.absolute, this[UID](entry), this[GID](entry));
                    } catch (er) {}
                }
                entry.resume();
            }
            [MKDIR](dir, mode) {
                try {
                    return mkdir.sync(dir, {
                        uid: this.uid,
                        gid: this.gid,
                        processUid: this.processUid,
                        processGid: this.processGid,
                        umask: this.processUmask,
                        preserve: this.preservePaths,
                        unlink: this.unlink,
                        cache: this.dirCache,
                        cwd: this.cwd,
                        mode
                    });
                } catch (er) {
                    return er;
                }
            }
            [LINK](entry, linkpath, link, _) {
                try {
                    fs[link + "Sync"](linkpath, entry.absolute);
                    entry.resume();
                } catch (er) {
                    return this[ONERROR](er, entry);
                }
            }
        }
        Unpack.Sync = UnpackSync;
        module.exports = Unpack;
    },
    4229: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const hlo = __webpack_require__(7461);
        const r = __webpack_require__(3666);
        module.exports = (opt_, files, cb) => {
            const opt = hlo(opt_);
            if (!opt.file) throw new TypeError("file is required");
            if (opt.gzip) throw new TypeError("cannot append to compressed archives");
            if (!files || !Array.isArray(files) || !files.length) throw new TypeError("no files or directories specified");
            files = Array.from(files);
            mtimeFilter(opt);
            return r(opt, files, cb);
        };
        const mtimeFilter = opt => {
            const filter = opt.filter;
            if (!opt.mtimeCache) opt.mtimeCache = new Map;
            opt.filter = filter ? (path, stat) => filter(path, stat) && !(opt.mtimeCache.get(path) > stat.mtime) : (path, stat) => !(opt.mtimeCache.get(path) > stat.mtime);
        };
    },
    8783: module => {
        "use strict";
        module.exports = Base => class extends Base {
            warn(code, message, data = {}) {
                if (this.file) data.file = this.file;
                if (this.cwd) data.cwd = this.cwd;
                data.code = message instanceof Error && message.code || code;
                data.tarCode = code;
                if (!this.strict && data.recoverable !== false) {
                    if (message instanceof Error) {
                        data = Object.assign(message, data);
                        message = message.message;
                    }
                    this.emit("warn", data.tarCode, message, data);
                } else if (message instanceof Error) this.emit("error", Object.assign(message, data)); else this.emit("error", Object.assign(new Error(`${code}: ${message}`), data));
            }
        };
    },
    6564: module => {
        "use strict";
        const raw = [ "|", "<", ">", "?", ":" ];
        const win = raw.map((char => String.fromCharCode(61440 + char.charCodeAt(0))));
        const toWin = new Map(raw.map(((char, i) => [ char, win[i] ])));
        const toRaw = new Map(win.map(((char, i) => [ char, raw[i] ])));
        module.exports = {
            encode: s => raw.reduce(((s, c) => s.split(c).join(toWin.get(c))), s),
            decode: s => win.reduce(((s, c) => s.split(c).join(toRaw.get(c))), s)
        };
    },
    8418: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        const MiniPass = __webpack_require__(2253);
        const Pax = __webpack_require__(9154);
        const Header = __webpack_require__(5017);
        const fs = __webpack_require__(5747);
        const path = __webpack_require__(5622);
        const maxReadSize = 16 * 1024 * 1024;
        const PROCESS = Symbol("process");
        const FILE = Symbol("file");
        const DIRECTORY = Symbol("directory");
        const SYMLINK = Symbol("symlink");
        const HARDLINK = Symbol("hardlink");
        const HEADER = Symbol("header");
        const READ = Symbol("read");
        const LSTAT = Symbol("lstat");
        const ONLSTAT = Symbol("onlstat");
        const ONREAD = Symbol("onread");
        const ONREADLINK = Symbol("onreadlink");
        const OPENFILE = Symbol("openfile");
        const ONOPENFILE = Symbol("onopenfile");
        const CLOSE = Symbol("close");
        const MODE = Symbol("mode");
        const warner = __webpack_require__(8783);
        const winchars = __webpack_require__(6564);
        const modeFix = __webpack_require__(9574);
        const WriteEntry = warner(class WriteEntry extends MiniPass {
            constructor(p, opt) {
                opt = opt || {};
                super(opt);
                if (typeof p !== "string") throw new TypeError("path is required");
                this.path = p;
                this.portable = !!opt.portable;
                this.myuid = process.getuid && process.getuid();
                this.myuser = process.env.USER || "";
                this.maxReadSize = opt.maxReadSize || maxReadSize;
                this.linkCache = opt.linkCache || new Map;
                this.statCache = opt.statCache || new Map;
                this.preservePaths = !!opt.preservePaths;
                this.cwd = opt.cwd || process.cwd();
                this.strict = !!opt.strict;
                this.noPax = !!opt.noPax;
                this.noMtime = !!opt.noMtime;
                this.mtime = opt.mtime || null;
                if (typeof opt.onwarn === "function") this.on("warn", opt.onwarn);
                let pathWarn = false;
                if (!this.preservePaths && path.win32.isAbsolute(p)) {
                    const parsed = path.win32.parse(p);
                    this.path = p.substr(parsed.root.length);
                    pathWarn = parsed.root;
                }
                this.win32 = !!opt.win32 || process.platform === "win32";
                if (this.win32) {
                    this.path = winchars.decode(this.path.replace(/\\/g, "/"));
                    p = p.replace(/\\/g, "/");
                }
                this.absolute = opt.absolute || path.resolve(this.cwd, p);
                if (this.path === "") this.path = "./";
                if (pathWarn) {
                    this.warn("TAR_ENTRY_INFO", `stripping ${pathWarn} from absolute path`, {
                        entry: this,
                        path: pathWarn + this.path
                    });
                }
                if (this.statCache.has(this.absolute)) this[ONLSTAT](this.statCache.get(this.absolute)); else this[LSTAT]();
            }
            [LSTAT]() {
                fs.lstat(this.absolute, ((er, stat) => {
                    if (er) return this.emit("error", er);
                    this[ONLSTAT](stat);
                }));
            }
            [ONLSTAT](stat) {
                this.statCache.set(this.absolute, stat);
                this.stat = stat;
                if (!stat.isFile()) stat.size = 0;
                this.type = getType(stat);
                this.emit("stat", stat);
                this[PROCESS]();
            }
            [PROCESS]() {
                switch (this.type) {
                  case "File":
                    return this[FILE]();

                  case "Directory":
                    return this[DIRECTORY]();

                  case "SymbolicLink":
                    return this[SYMLINK]();

                  default:
                    return this.end();
                }
            }
            [MODE](mode) {
                return modeFix(mode, this.type === "Directory", this.portable);
            }
            [HEADER]() {
                if (this.type === "Directory" && this.portable) this.noMtime = true;
                this.header = new Header({
                    path: this.path,
                    linkpath: this.linkpath,
                    mode: this[MODE](this.stat.mode),
                    uid: this.portable ? null : this.stat.uid,
                    gid: this.portable ? null : this.stat.gid,
                    size: this.stat.size,
                    mtime: this.noMtime ? null : this.mtime || this.stat.mtime,
                    type: this.type,
                    uname: this.portable ? null : this.stat.uid === this.myuid ? this.myuser : "",
                    atime: this.portable ? null : this.stat.atime,
                    ctime: this.portable ? null : this.stat.ctime
                });
                if (this.header.encode() && !this.noPax) {
                    this.write(new Pax({
                        atime: this.portable ? null : this.header.atime,
                        ctime: this.portable ? null : this.header.ctime,
                        gid: this.portable ? null : this.header.gid,
                        mtime: this.noMtime ? null : this.mtime || this.header.mtime,
                        path: this.path,
                        linkpath: this.linkpath,
                        size: this.header.size,
                        uid: this.portable ? null : this.header.uid,
                        uname: this.portable ? null : this.header.uname,
                        dev: this.portable ? null : this.stat.dev,
                        ino: this.portable ? null : this.stat.ino,
                        nlink: this.portable ? null : this.stat.nlink
                    }).encode());
                }
                this.write(this.header.block);
            }
            [DIRECTORY]() {
                if (this.path.substr(-1) !== "/") this.path += "/";
                this.stat.size = 0;
                this[HEADER]();
                this.end();
            }
            [SYMLINK]() {
                fs.readlink(this.absolute, ((er, linkpath) => {
                    if (er) return this.emit("error", er);
                    this[ONREADLINK](linkpath);
                }));
            }
            [ONREADLINK](linkpath) {
                this.linkpath = linkpath.replace(/\\/g, "/");
                this[HEADER]();
                this.end();
            }
            [HARDLINK](linkpath) {
                this.type = "Link";
                this.linkpath = path.relative(this.cwd, linkpath).replace(/\\/g, "/");
                this.stat.size = 0;
                this[HEADER]();
                this.end();
            }
            [FILE]() {
                if (this.stat.nlink > 1) {
                    const linkKey = this.stat.dev + ":" + this.stat.ino;
                    if (this.linkCache.has(linkKey)) {
                        const linkpath = this.linkCache.get(linkKey);
                        if (linkpath.indexOf(this.cwd) === 0) return this[HARDLINK](linkpath);
                    }
                    this.linkCache.set(linkKey, this.absolute);
                }
                this[HEADER]();
                if (this.stat.size === 0) return this.end();
                this[OPENFILE]();
            }
            [OPENFILE]() {
                fs.open(this.absolute, "r", ((er, fd) => {
                    if (er) return this.emit("error", er);
                    this[ONOPENFILE](fd);
                }));
            }
            [ONOPENFILE](fd) {
                const blockLen = 512 * Math.ceil(this.stat.size / 512);
                const bufLen = Math.min(blockLen, this.maxReadSize);
                const buf = Buffer.allocUnsafe(bufLen);
                this[READ](fd, buf, 0, buf.length, 0, this.stat.size, blockLen);
            }
            [READ](fd, buf, offset, length, pos, remain, blockRemain) {
                fs.read(fd, buf, offset, length, pos, ((er, bytesRead) => {
                    if (er) {
                        return this[CLOSE](fd, (() => this.emit("error", er)));
                    }
                    this[ONREAD](fd, buf, offset, length, pos, remain, blockRemain, bytesRead);
                }));
            }
            [CLOSE](fd, cb) {
                fs.close(fd, cb);
            }
            [ONREAD](fd, buf, offset, length, pos, remain, blockRemain, bytesRead) {
                if (bytesRead <= 0 && remain > 0) {
                    const er = new Error("encountered unexpected EOF");
                    er.path = this.absolute;
                    er.syscall = "read";
                    er.code = "EOF";
                    return this[CLOSE](fd, (() => this.emit("error", er)));
                }
                if (bytesRead > remain) {
                    const er = new Error("did not encounter expected EOF");
                    er.path = this.absolute;
                    er.syscall = "read";
                    er.code = "EOF";
                    return this[CLOSE](fd, (() => this.emit("error", er)));
                }
                if (bytesRead === remain) {
                    for (let i = bytesRead; i < length && bytesRead < blockRemain; i++) {
                        buf[i + offset] = 0;
                        bytesRead++;
                        remain++;
                    }
                }
                const writeBuf = offset === 0 && bytesRead === buf.length ? buf : buf.slice(offset, offset + bytesRead);
                remain -= bytesRead;
                blockRemain -= bytesRead;
                pos += bytesRead;
                offset += bytesRead;
                this.write(writeBuf);
                if (!remain) {
                    if (blockRemain) this.write(Buffer.alloc(blockRemain));
                    return this[CLOSE](fd, (er => er ? this.emit("error", er) : this.end()));
                }
                if (offset >= length) {
                    buf = Buffer.allocUnsafe(length);
                    offset = 0;
                }
                length = buf.length - offset;
                this[READ](fd, buf, offset, length, pos, remain, blockRemain);
            }
        });
        class WriteEntrySync extends WriteEntry {
            [LSTAT]() {
                this[ONLSTAT](fs.lstatSync(this.absolute));
            }
            [SYMLINK]() {
                this[ONREADLINK](fs.readlinkSync(this.absolute));
            }
            [OPENFILE]() {
                this[ONOPENFILE](fs.openSync(this.absolute, "r"));
            }
            [READ](fd, buf, offset, length, pos, remain, blockRemain) {
                let threw = true;
                try {
                    const bytesRead = fs.readSync(fd, buf, offset, length, pos);
                    this[ONREAD](fd, buf, offset, length, pos, remain, blockRemain, bytesRead);
                    threw = false;
                } finally {
                    if (threw) {
                        try {
                            this[CLOSE](fd, (() => {}));
                        } catch (er) {}
                    }
                }
            }
            [CLOSE](fd, cb) {
                fs.closeSync(fd);
                cb();
            }
        }
        const WriteEntryTar = warner(class WriteEntryTar extends MiniPass {
            constructor(readEntry, opt) {
                opt = opt || {};
                super(opt);
                this.preservePaths = !!opt.preservePaths;
                this.portable = !!opt.portable;
                this.strict = !!opt.strict;
                this.noPax = !!opt.noPax;
                this.noMtime = !!opt.noMtime;
                this.readEntry = readEntry;
                this.type = readEntry.type;
                if (this.type === "Directory" && this.portable) this.noMtime = true;
                this.path = readEntry.path;
                this.mode = this[MODE](readEntry.mode);
                this.uid = this.portable ? null : readEntry.uid;
                this.gid = this.portable ? null : readEntry.gid;
                this.uname = this.portable ? null : readEntry.uname;
                this.gname = this.portable ? null : readEntry.gname;
                this.size = readEntry.size;
                this.mtime = this.noMtime ? null : opt.mtime || readEntry.mtime;
                this.atime = this.portable ? null : readEntry.atime;
                this.ctime = this.portable ? null : readEntry.ctime;
                this.linkpath = readEntry.linkpath;
                if (typeof opt.onwarn === "function") this.on("warn", opt.onwarn);
                let pathWarn = false;
                if (path.isAbsolute(this.path) && !this.preservePaths) {
                    const parsed = path.parse(this.path);
                    pathWarn = parsed.root;
                    this.path = this.path.substr(parsed.root.length);
                }
                this.remain = readEntry.size;
                this.blockRemain = readEntry.startBlockSize;
                this.header = new Header({
                    path: this.path,
                    linkpath: this.linkpath,
                    mode: this.mode,
                    uid: this.portable ? null : this.uid,
                    gid: this.portable ? null : this.gid,
                    size: this.size,
                    mtime: this.noMtime ? null : this.mtime,
                    type: this.type,
                    uname: this.portable ? null : this.uname,
                    atime: this.portable ? null : this.atime,
                    ctime: this.portable ? null : this.ctime
                });
                if (pathWarn) {
                    this.warn("TAR_ENTRY_INFO", `stripping ${pathWarn} from absolute path`, {
                        entry: this,
                        path: pathWarn + this.path
                    });
                }
                if (this.header.encode() && !this.noPax) {
                    super.write(new Pax({
                        atime: this.portable ? null : this.atime,
                        ctime: this.portable ? null : this.ctime,
                        gid: this.portable ? null : this.gid,
                        mtime: this.noMtime ? null : this.mtime,
                        path: this.path,
                        linkpath: this.linkpath,
                        size: this.size,
                        uid: this.portable ? null : this.uid,
                        uname: this.portable ? null : this.uname,
                        dev: this.portable ? null : this.readEntry.dev,
                        ino: this.portable ? null : this.readEntry.ino,
                        nlink: this.portable ? null : this.readEntry.nlink
                    }).encode());
                }
                super.write(this.header.block);
                readEntry.pipe(this);
            }
            [MODE](mode) {
                return modeFix(mode, this.type === "Directory", this.portable);
            }
            write(data) {
                const writeLen = data.length;
                if (writeLen > this.blockRemain) throw new Error("writing more to entry than is appropriate");
                this.blockRemain -= writeLen;
                return super.write(data);
            }
            end() {
                if (this.blockRemain) this.write(Buffer.alloc(this.blockRemain));
                return super.end();
            }
        });
        WriteEntry.Sync = WriteEntrySync;
        WriteEntry.Tar = WriteEntryTar;
        const getType = stat => stat.isFile() ? "File" : stat.isDirectory() ? "Directory" : stat.isSymbolicLink() ? "SymbolicLink" : "Unsupported";
        module.exports = WriteEntry;
    },
    3278: module => {
        "use strict";
        module.exports = function(Yallist) {
            Yallist.prototype[Symbol.iterator] = function*() {
                for (let walker = this.head; walker; walker = walker.next) {
                    yield walker.value;
                }
            };
        };
    },
    1455: (module, __unused_webpack_exports, __webpack_require__) => {
        "use strict";
        module.exports = Yallist;
        Yallist.Node = Node;
        Yallist.create = Yallist;
        function Yallist(list) {
            var self = this;
            if (!(self instanceof Yallist)) {
                self = new Yallist;
            }
            self.tail = null;
            self.head = null;
            self.length = 0;
            if (list && typeof list.forEach === "function") {
                list.forEach((function(item) {
                    self.push(item);
                }));
            } else if (arguments.length > 0) {
                for (var i = 0, l = arguments.length; i < l; i++) {
                    self.push(arguments[i]);
                }
            }
            return self;
        }
        Yallist.prototype.removeNode = function(node) {
            if (node.list !== this) {
                throw new Error("removing node which does not belong to this list");
            }
            var next = node.next;
            var prev = node.prev;
            if (next) {
                next.prev = prev;
            }
            if (prev) {
                prev.next = next;
            }
            if (node === this.head) {
                this.head = next;
            }
            if (node === this.tail) {
                this.tail = prev;
            }
            node.list.length--;
            node.next = null;
            node.prev = null;
            node.list = null;
            return next;
        };
        Yallist.prototype.unshiftNode = function(node) {
            if (node === this.head) {
                return;
            }
            if (node.list) {
                node.list.removeNode(node);
            }
            var head = this.head;
            node.list = this;
            node.next = head;
            if (head) {
                head.prev = node;
            }
            this.head = node;
            if (!this.tail) {
                this.tail = node;
            }
            this.length++;
        };
        Yallist.prototype.pushNode = function(node) {
            if (node === this.tail) {
                return;
            }
            if (node.list) {
                node.list.removeNode(node);
            }
            var tail = this.tail;
            node.list = this;
            node.prev = tail;
            if (tail) {
                tail.next = node;
            }
            this.tail = node;
            if (!this.head) {
                this.head = node;
            }
            this.length++;
        };
        Yallist.prototype.push = function() {
            for (var i = 0, l = arguments.length; i < l; i++) {
                push(this, arguments[i]);
            }
            return this.length;
        };
        Yallist.prototype.unshift = function() {
            for (var i = 0, l = arguments.length; i < l; i++) {
                unshift(this, arguments[i]);
            }
            return this.length;
        };
        Yallist.prototype.pop = function() {
            if (!this.tail) {
                return undefined;
            }
            var res = this.tail.value;
            this.tail = this.tail.prev;
            if (this.tail) {
                this.tail.next = null;
            } else {
                this.head = null;
            }
            this.length--;
            return res;
        };
        Yallist.prototype.shift = function() {
            if (!this.head) {
                return undefined;
            }
            var res = this.head.value;
            this.head = this.head.next;
            if (this.head) {
                this.head.prev = null;
            } else {
                this.tail = null;
            }
            this.length--;
            return res;
        };
        Yallist.prototype.forEach = function(fn, thisp) {
            thisp = thisp || this;
            for (var walker = this.head, i = 0; walker !== null; i++) {
                fn.call(thisp, walker.value, i, this);
                walker = walker.next;
            }
        };
        Yallist.prototype.forEachReverse = function(fn, thisp) {
            thisp = thisp || this;
            for (var walker = this.tail, i = this.length - 1; walker !== null; i--) {
                fn.call(thisp, walker.value, i, this);
                walker = walker.prev;
            }
        };
        Yallist.prototype.get = function(n) {
            for (var i = 0, walker = this.head; walker !== null && i < n; i++) {
                walker = walker.next;
            }
            if (i === n && walker !== null) {
                return walker.value;
            }
        };
        Yallist.prototype.getReverse = function(n) {
            for (var i = 0, walker = this.tail; walker !== null && i < n; i++) {
                walker = walker.prev;
            }
            if (i === n && walker !== null) {
                return walker.value;
            }
        };
        Yallist.prototype.map = function(fn, thisp) {
            thisp = thisp || this;
            var res = new Yallist;
            for (var walker = this.head; walker !== null; ) {
                res.push(fn.call(thisp, walker.value, this));
                walker = walker.next;
            }
            return res;
        };
        Yallist.prototype.mapReverse = function(fn, thisp) {
            thisp = thisp || this;
            var res = new Yallist;
            for (var walker = this.tail; walker !== null; ) {
                res.push(fn.call(thisp, walker.value, this));
                walker = walker.prev;
            }
            return res;
        };
        Yallist.prototype.reduce = function(fn, initial) {
            var acc;
            var walker = this.head;
            if (arguments.length > 1) {
                acc = initial;
            } else if (this.head) {
                walker = this.head.next;
                acc = this.head.value;
            } else {
                throw new TypeError("Reduce of empty list with no initial value");
            }
            for (var i = 0; walker !== null; i++) {
                acc = fn(acc, walker.value, i);
                walker = walker.next;
            }
            return acc;
        };
        Yallist.prototype.reduceReverse = function(fn, initial) {
            var acc;
            var walker = this.tail;
            if (arguments.length > 1) {
                acc = initial;
            } else if (this.tail) {
                walker = this.tail.prev;
                acc = this.tail.value;
            } else {
                throw new TypeError("Reduce of empty list with no initial value");
            }
            for (var i = this.length - 1; walker !== null; i--) {
                acc = fn(acc, walker.value, i);
                walker = walker.prev;
            }
            return acc;
        };
        Yallist.prototype.toArray = function() {
            var arr = new Array(this.length);
            for (var i = 0, walker = this.head; walker !== null; i++) {
                arr[i] = walker.value;
                walker = walker.next;
            }
            return arr;
        };
        Yallist.prototype.toArrayReverse = function() {
            var arr = new Array(this.length);
            for (var i = 0, walker = this.tail; walker !== null; i++) {
                arr[i] = walker.value;
                walker = walker.prev;
            }
            return arr;
        };
        Yallist.prototype.slice = function(from, to) {
            to = to || this.length;
            if (to < 0) {
                to += this.length;
            }
            from = from || 0;
            if (from < 0) {
                from += this.length;
            }
            var ret = new Yallist;
            if (to < from || to < 0) {
                return ret;
            }
            if (from < 0) {
                from = 0;
            }
            if (to > this.length) {
                to = this.length;
            }
            for (var i = 0, walker = this.head; walker !== null && i < from; i++) {
                walker = walker.next;
            }
            for (;walker !== null && i < to; i++, walker = walker.next) {
                ret.push(walker.value);
            }
            return ret;
        };
        Yallist.prototype.sliceReverse = function(from, to) {
            to = to || this.length;
            if (to < 0) {
                to += this.length;
            }
            from = from || 0;
            if (from < 0) {
                from += this.length;
            }
            var ret = new Yallist;
            if (to < from || to < 0) {
                return ret;
            }
            if (from < 0) {
                from = 0;
            }
            if (to > this.length) {
                to = this.length;
            }
            for (var i = this.length, walker = this.tail; walker !== null && i > to; i--) {
                walker = walker.prev;
            }
            for (;walker !== null && i > from; i--, walker = walker.prev) {
                ret.push(walker.value);
            }
            return ret;
        };
        Yallist.prototype.splice = function(start, deleteCount, ...nodes) {
            if (start > this.length) {
                start = this.length - 1;
            }
            if (start < 0) {
                start = this.length + start;
            }
            for (var i = 0, walker = this.head; walker !== null && i < start; i++) {
                walker = walker.next;
            }
            var ret = [];
            for (var i = 0; walker && i < deleteCount; i++) {
                ret.push(walker.value);
                walker = this.removeNode(walker);
            }
            if (walker === null) {
                walker = this.tail;
            }
            if (walker !== this.head && walker !== this.tail) {
                walker = walker.prev;
            }
            for (var i = 0; i < nodes.length; i++) {
                walker = insert(this, walker, nodes[i]);
            }
            return ret;
        };
        Yallist.prototype.reverse = function() {
            var head = this.head;
            var tail = this.tail;
            for (var walker = head; walker !== null; walker = walker.prev) {
                var p = walker.prev;
                walker.prev = walker.next;
                walker.next = p;
            }
            this.head = tail;
            this.tail = head;
            return this;
        };
        function insert(self, node, value) {
            var inserted = node === self.head ? new Node(value, null, node, self) : new Node(value, node, node.next, self);
            if (inserted.next === null) {
                self.tail = inserted;
            }
            if (inserted.prev === null) {
                self.head = inserted;
            }
            self.length++;
            return inserted;
        }
        function push(self, item) {
            self.tail = new Node(item, self.tail, null, self);
            if (!self.head) {
                self.head = self.tail;
            }
            self.length++;
        }
        function unshift(self, item) {
            self.head = new Node(item, null, self.head, self);
            if (!self.tail) {
                self.tail = self.head;
            }
            self.length++;
        }
        function Node(value, prev, next, list) {
            if (!(this instanceof Node)) {
                return new Node(value, prev, next, list);
            }
            this.list = list;
            this.value = value;
            if (prev) {
                prev.next = this;
                this.prev = prev;
            } else {
                this.prev = null;
            }
            if (next) {
                next.prev = this;
                this.next = next;
            } else {
                this.next = null;
            }
        }
        try {
            __webpack_require__(3278)(Yallist);
        } catch (er) {}
    },
    2816: (__unused_webpack_module, exports) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.isPropertyOverride = exports.isMethodOverride = exports.isWireStruct = exports.isWireMap = exports.isWireEnum = exports.isWireDate = exports.isObjRef = exports.TOKEN_STRUCT = exports.TOKEN_MAP = exports.TOKEN_ENUM = exports.TOKEN_DATE = exports.TOKEN_INTERFACES = exports.TOKEN_REF = void 0;
        exports.TOKEN_REF = "$jsii.byref";
        exports.TOKEN_INTERFACES = "$jsii.interfaces";
        exports.TOKEN_DATE = "$jsii.date";
        exports.TOKEN_ENUM = "$jsii.enum";
        exports.TOKEN_MAP = "$jsii.map";
        exports.TOKEN_STRUCT = "$jsii.struct";
        function isObjRef(value) {
            return typeof value === "object" && value !== null && exports.TOKEN_REF in value;
        }
        exports.isObjRef = isObjRef;
        function isWireDate(value) {
            return typeof value === "object" && value !== null && exports.TOKEN_DATE in value;
        }
        exports.isWireDate = isWireDate;
        function isWireEnum(value) {
            return typeof value === "object" && value !== null && exports.TOKEN_ENUM in value;
        }
        exports.isWireEnum = isWireEnum;
        function isWireMap(value) {
            return typeof value === "object" && value !== null && exports.TOKEN_MAP in value;
        }
        exports.isWireMap = isWireMap;
        function isWireStruct(value) {
            return typeof value === "object" && value !== null && exports.TOKEN_STRUCT in value;
        }
        exports.isWireStruct = isWireStruct;
        function isMethodOverride(value) {
            return value.method != null;
        }
        exports.isMethodOverride = isMethodOverride;
        function isPropertyOverride(value) {
            return value.property != null;
        }
        exports.isPropertyOverride = isPropertyOverride;
    },
    8944: function(__unused_webpack_module, exports, __webpack_require__) {
        "use strict";
        var __createBinding = this && this.__createBinding || (Object.create ? function(o, m, k, k2) {
            if (k2 === undefined) k2 = k;
            Object.defineProperty(o, k2, {
                enumerable: true,
                get: function() {
                    return m[k];
                }
            });
        } : function(o, m, k, k2) {
            if (k2 === undefined) k2 = k;
            o[k2] = m[k];
        });
        var __exportStar = this && this.__exportStar || function(m, exports) {
            for (var p in m) if (p !== "default" && !exports.hasOwnProperty(p)) __createBinding(exports, m, p);
        };
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.api = void 0;
        __exportStar(__webpack_require__(2742), exports);
        const api = __webpack_require__(2816);
        exports.api = api;
    },
    2742: (module, exports, __webpack_require__) => {
        "use strict";
        module = __webpack_require__.nmd(module);
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.Kernel = void 0;
        const spec = __webpack_require__(1804);
        const cp = __webpack_require__(3129);
        const fs = __webpack_require__(9728);
        const os = __webpack_require__(2087);
        const path = __webpack_require__(5622);
        const tar = __webpack_require__(1189);
        const vm = __webpack_require__(2184);
        const api = __webpack_require__(2816);
        const api_1 = __webpack_require__(2816);
        const objects_1 = __webpack_require__(2309);
        const onExit = __webpack_require__(6703);
        const wire = __webpack_require__(3837);
        class Kernel {
            constructor(callbackHandler) {
                this.callbackHandler = callbackHandler;
                this.traceEnabled = false;
                this.assemblies = {};
                this.objects = new objects_1.ObjectTable(this._typeInfoForFqn.bind(this));
                this.cbs = {};
                this.waiting = {};
                this.promises = {};
                this.nextid = 2e4;
                const moduleLoad = __webpack_require__(2282).Module._load;
                const nodeRequire = p => moduleLoad(p, module, false);
                this.sandbox = vm.createContext({
                    Buffer,
                    setImmediate,
                    require: nodeRequire
                });
            }
            load(req) {
                var _a, _b;
                this._debug("load", req);
                if ("assembly" in req) {
                    throw new Error('`assembly` field is deprecated for "load", use `name`, `version` and `tarball` instead');
                }
                const pkgname = req.name;
                const pkgver = req.version;
                const packageDir = this._getPackageDir(pkgname);
                if (fs.pathExistsSync(packageDir)) {
                    const epkg = fs.readJsonSync(path.join(packageDir, "package.json"));
                    if (epkg.version !== pkgver) {
                        throw new Error(`Multiple versions ${pkgver} and ${epkg.version} of the ` + `package '${pkgname}' cannot be loaded together since this is unsupported by ` + "some runtime environments");
                    }
                    this._debug("look up already-loaded assembly", pkgname);
                    const assm = this.assemblies[pkgname];
                    return {
                        assembly: assm.metadata.name,
                        types: Object.keys((_a = assm.metadata.types) !== null && _a !== void 0 ? _a : {}).length
                    };
                }
                fs.mkdirpSync(packageDir);
                const originalUmask = process.umask(18);
                try {
                    tar.extract({
                        cwd: packageDir,
                        file: req.tarball,
                        strict: true,
                        strip: 1,
                        sync: true,
                        unlink: true
                    });
                } finally {
                    process.umask(originalUmask);
                }
                const jsiiMetadataFile = path.join(packageDir, spec.SPEC_FILE_NAME);
                if (!fs.pathExistsSync(jsiiMetadataFile)) {
                    throw new Error(`Package tarball ${req.tarball} must have a file named ${spec.SPEC_FILE_NAME} at the root`);
                }
                const assmSpec = fs.readJsonSync(jsiiMetadataFile);
                const closure = this._execute(`require(String.raw\`${packageDir}\`)`, packageDir);
                const assm = new Assembly(assmSpec, closure);
                this._addAssembly(assm);
                return {
                    assembly: assmSpec.name,
                    types: Object.keys((_b = assmSpec.types) !== null && _b !== void 0 ? _b : {}).length
                };
            }
            invokeBinScript(req) {
                var _a;
                const packageDir = this._getPackageDir(req.assembly);
                if (fs.pathExistsSync(packageDir)) {
                    const epkg = fs.readJsonSync(path.join(packageDir, "package.json"));
                    if (!epkg.bin) {
                        throw new Error("There is no bin scripts defined for this package.");
                    }
                    const scriptPath = epkg.bin[req.script];
                    if (!epkg.bin) {
                        throw new Error(`Script with name ${req.script} was not defined.`);
                    }
                    const result = cp.spawnSync(path.join(packageDir, scriptPath), (_a = req.args) !== null && _a !== void 0 ? _a : [], {
                        encoding: "utf-8",
                        env: {
                            ...process.env,
                            NODE_OPTIONS: process.execArgv.join(" "),
                            PATH: `${path.dirname(process.execPath)}:${process.env.PATH}`
                        },
                        shell: true
                    });
                    return {
                        stdout: result.stdout,
                        stderr: result.stderr,
                        status: result.status,
                        signal: result.signal
                    };
                }
                throw new Error(`Package with name ${req.assembly} was not loaded.`);
            }
            create(req) {
                return this._create(req);
            }
            del(req) {
                const {objref} = req;
                this._debug("del", objref);
                this.objects.deleteObject(objref);
                return {};
            }
            sget(req) {
                const {fqn, property} = req;
                const symbol = `${fqn}.${property}`;
                this._debug("sget", symbol);
                const ti = this._typeInfoForProperty(property, fqn);
                if (!ti.static) {
                    throw new Error(`property ${symbol} is not static`);
                }
                const prototype = this._findSymbol(fqn);
                const value = this._ensureSync(`property ${property}`, (() => this._wrapSandboxCode((() => prototype[property]))));
                this._debug("value:", value);
                const ret = this._fromSandbox(value, ti);
                this._debug("ret", ret);
                return {
                    value: ret
                };
            }
            sset(req) {
                const {fqn, property, value} = req;
                const symbol = `${fqn}.${property}`;
                this._debug("sset", symbol);
                const ti = this._typeInfoForProperty(property, fqn);
                if (!ti.static) {
                    throw new Error(`property ${symbol} is not static`);
                }
                if (ti.immutable) {
                    throw new Error(`static property ${symbol} is readonly`);
                }
                const prototype = this._findSymbol(fqn);
                this._ensureSync(`property ${property}`, (() => this._wrapSandboxCode((() => prototype[property] = this._toSandbox(value, ti)))));
                return {};
            }
            get(req) {
                const {objref, property} = req;
                this._debug("get", objref, property);
                const {instance, fqn, interfaces} = this.objects.findObject(objref);
                const ti = this._typeInfoForProperty(property, fqn, interfaces);
                const propertyToGet = this._findPropertyTarget(instance, property);
                const value = this._ensureSync(`property '${objref[api_1.TOKEN_REF]}.${propertyToGet}'`, (() => this._wrapSandboxCode((() => instance[propertyToGet]))));
                this._debug("value:", value);
                const ret = this._fromSandbox(value, ti);
                this._debug("ret:", ret);
                return {
                    value: ret
                };
            }
            set(req) {
                const {objref, property, value} = req;
                this._debug("set", objref, property, value);
                const {instance, fqn, interfaces} = this.objects.findObject(objref);
                const propInfo = this._typeInfoForProperty(req.property, fqn, interfaces);
                if (propInfo.immutable) {
                    throw new Error(`Cannot set value of immutable property ${req.property} to ${req.value}`);
                }
                const propertyToSet = this._findPropertyTarget(instance, property);
                this._ensureSync(`property '${objref[api_1.TOKEN_REF]}.${propertyToSet}'`, (() => this._wrapSandboxCode((() => instance[propertyToSet] = this._toSandbox(value, propInfo)))));
                return {};
            }
            invoke(req) {
                var _a, _b;
                const {objref, method} = req;
                const args = (_a = req.args) !== null && _a !== void 0 ? _a : [];
                this._debug("invoke", objref, method, args);
                const {ti, obj, fn} = this._findInvokeTarget(objref, method, args);
                if (ti.async) {
                    throw new Error(`${method} is an async method, use "begin" instead`);
                }
                const ret = this._ensureSync(`method '${objref[api_1.TOKEN_REF]}.${method}'`, (() => this._wrapSandboxCode((() => fn.apply(obj, this._toSandboxValues(args, ti.parameters))))));
                const result = this._fromSandbox(ret, (_b = ti.returns) !== null && _b !== void 0 ? _b : "void");
                this._debug("invoke result", result);
                return {
                    result
                };
            }
            sinvoke(req) {
                var _a, _b;
                const {fqn, method} = req;
                const args = (_a = req.args) !== null && _a !== void 0 ? _a : [];
                this._debug("sinvoke", fqn, method, args);
                const ti = this._typeInfoForMethod(method, fqn);
                if (!ti.static) {
                    throw new Error(`${fqn}.${method} is not a static method`);
                }
                if (ti.async) {
                    throw new Error(`${method} is an async method, use "begin" instead`);
                }
                const prototype = this._findSymbol(fqn);
                const fn = prototype[method];
                const ret = this._ensureSync(`method '${fqn}.${method}'`, (() => this._wrapSandboxCode((() => fn.apply(prototype, this._toSandboxValues(args, ti.parameters))))));
                this._debug("method returned:", ret);
                return {
                    result: this._fromSandbox(ret, (_b = ti.returns) !== null && _b !== void 0 ? _b : "void")
                };
            }
            begin(req) {
                var _a;
                const {objref, method} = req;
                const args = (_a = req.args) !== null && _a !== void 0 ? _a : [];
                this._debug("begin", objref, method, args);
                if (this.syncInProgress) {
                    throw new Error(`Cannot invoke async method '${req.objref[api_1.TOKEN_REF]}.${req.method}' while sync ${this.syncInProgress} is being processed`);
                }
                const {ti, obj, fn} = this._findInvokeTarget(objref, method, args);
                if (!ti.async) {
                    throw new Error(`Method ${method} is expected to be an async method`);
                }
                const promise = this._wrapSandboxCode((() => fn.apply(obj, this._toSandboxValues(args, ti.parameters))));
                promise.catch((_ => undefined));
                const prid = this._makeprid();
                this.promises[prid] = {
                    promise,
                    method: ti
                };
                return {
                    promiseid: prid
                };
            }
            async end(req) {
                var _a;
                const {promiseid} = req;
                this._debug("end", promiseid);
                const {promise, method} = this.promises[promiseid];
                if (promise == null) {
                    throw new Error(`Cannot find promise with ID: ${promiseid}`);
                }
                let result;
                try {
                    result = await promise;
                    this._debug("promise result:", result);
                } catch (e) {
                    this._debug("promise error:", e);
                    throw e;
                }
                return {
                    result: this._fromSandbox(result, (_a = method.returns) !== null && _a !== void 0 ? _a : "void")
                };
            }
            callbacks(_req) {
                this._debug("callbacks");
                const ret = Object.keys(this.cbs).map((cbid => {
                    const cb = this.cbs[cbid];
                    this.waiting[cbid] = cb;
                    const callback = {
                        cbid,
                        cookie: cb.override.cookie,
                        invoke: {
                            objref: cb.objref,
                            method: cb.override.method,
                            args: cb.args
                        }
                    };
                    return callback;
                }));
                this.cbs = {};
                return {
                    callbacks: ret
                };
            }
            complete(req) {
                var _a;
                const {cbid, err, result} = req;
                this._debug("complete", cbid, err, result);
                if (!(cbid in this.waiting)) {
                    throw new Error(`Callback ${cbid} not found`);
                }
                const cb = this.waiting[cbid];
                if (err) {
                    this._debug("completed with error:", err);
                    cb.fail(new Error(err));
                } else {
                    const sandoxResult = this._toSandbox(result, (_a = cb.expectedReturnType) !== null && _a !== void 0 ? _a : "void");
                    this._debug("completed with result:", sandoxResult);
                    cb.succeed(sandoxResult);
                }
                delete this.waiting[cbid];
                return {
                    cbid
                };
            }
            naming(req) {
                const assemblyName = req.assembly;
                this._debug("naming", assemblyName);
                const assembly = this._assemblyFor(assemblyName);
                const targets = assembly.metadata.targets;
                if (!targets) {
                    throw new Error(`Unexpected - "targets" for ${assemblyName} is missing!`);
                }
                return {
                    naming: targets
                };
            }
            stats(_req) {
                return {
                    objectCount: this.objects.count
                };
            }
            _addAssembly(assm) {
                var _a;
                this.assemblies[assm.metadata.name] = assm;
                for (const fqn of Object.keys((_a = assm.metadata.types) !== null && _a !== void 0 ? _a : {})) {
                    const typedef = assm.metadata.types[fqn];
                    switch (typedef.kind) {
                      case spec.TypeKind.Interface:
                        continue;

                      case spec.TypeKind.Class:
                      case spec.TypeKind.Enum:
                        const constructor = this._findSymbol(fqn);
                        objects_1.tagJsiiConstructor(constructor, fqn);
                    }
                }
            }
            _findCtor(fqn, args) {
                if (fqn === wire.EMPTY_OBJECT_FQN) {
                    return {
                        ctor: Object
                    };
                }
                const typeinfo = this._typeInfoForFqn(fqn);
                switch (typeinfo.kind) {
                  case spec.TypeKind.Class:
                    const classType = typeinfo;
                    this._validateMethodArguments(classType.initializer, args);
                    return {
                        ctor: this._findSymbol(fqn),
                        parameters: classType.initializer && classType.initializer.parameters
                    };

                  case spec.TypeKind.Interface:
                    throw new Error(`Cannot create an object with an FQN of an interface: ${fqn}`);

                  default:
                    throw new Error(`Unexpected FQN kind: ${fqn}`);
                }
            }
            _getPackageDir(pkgname) {
                if (!this.installDir) {
                    this.installDir = fs.mkdtempSync(path.join(os.tmpdir(), "jsii-kernel-"));
                    fs.mkdirpSync(path.join(this.installDir, "node_modules"));
                    this._debug("creating jsii-kernel modules workdir:", this.installDir);
                    onExit.removeSync(this.installDir);
                }
                return path.join(this.installDir, "node_modules", pkgname);
            }
            _create(req) {
                var _a, _b;
                this._debug("create", req);
                const {fqn, interfaces, overrides} = req;
                const requestArgs = (_a = req.args) !== null && _a !== void 0 ? _a : [];
                const ctorResult = this._findCtor(fqn, requestArgs);
                const ctor = ctorResult.ctor;
                const obj = this._wrapSandboxCode((() => new ctor(...this._toSandboxValues(requestArgs, ctorResult.parameters))));
                const objref = this.objects.registerObject(obj, fqn, (_b = req.interfaces) !== null && _b !== void 0 ? _b : []);
                if (overrides) {
                    this._debug("overrides", overrides);
                    const overrideTypeErrorMessage = 'Override can either be "method" or "property"';
                    const methods = new Set;
                    const properties = new Set;
                    for (const override of overrides) {
                        if (api.isMethodOverride(override)) {
                            if (api.isPropertyOverride(override)) {
                                throw new Error(overrideTypeErrorMessage);
                            }
                            if (methods.has(override.method)) {
                                throw new Error(`Duplicate override for method '${override.method}'`);
                            }
                            methods.add(override.method);
                            this._applyMethodOverride(obj, objref, fqn, interfaces, override);
                        } else if (api.isPropertyOverride(override)) {
                            if (api.isMethodOverride(override)) {
                                throw new Error(overrideTypeErrorMessage);
                            }
                            if (properties.has(override.property)) {
                                throw Error(`Duplicate override for property '${override.property}'`);
                            }
                            properties.add(override.property);
                            this._applyPropertyOverride(obj, objref, fqn, interfaces, override);
                        } else {
                            throw new Error(overrideTypeErrorMessage);
                        }
                    }
                }
                return objref;
            }
            _getSuperPropertyName(name) {
                return `$jsii$super$${name}$`;
            }
            _applyPropertyOverride(obj, objref, typeFqn, interfaces, override) {
                if (this._tryTypeInfoForMethod(override.property, typeFqn, interfaces)) {
                    throw new Error(`Trying to override method '${override.property}' as a property`);
                }
                let propInfo = this._tryTypeInfoForProperty(override.property, typeFqn, interfaces);
                if (!propInfo && override.property in obj) {
                    this._debug(`Skipping override of private property ${override.property}`);
                    return;
                }
                if (!propInfo) {
                    propInfo = {
                        name: override.property,
                        type: spec.CANONICAL_ANY
                    };
                }
                this._defineOverridenProperty(obj, objref, override, propInfo);
            }
            _defineOverridenProperty(obj, objref, override, propInfo) {
                var _a;
                const propertyName = override.property;
                this._debug("apply override", propertyName);
                const prev = (_a = getPropertyDescriptor(obj, propertyName)) !== null && _a !== void 0 ? _a : {
                    value: obj[propertyName],
                    writable: true,
                    enumerable: true,
                    configurable: true
                };
                const prevEnumerable = prev.enumerable;
                prev.enumerable = false;
                Object.defineProperty(obj, this._getSuperPropertyName(propertyName), prev);
                Object.defineProperty(obj, propertyName, {
                    enumerable: prevEnumerable,
                    configurable: prev.configurable,
                    get: () => {
                        this._debug("virtual get", objref, propertyName, {
                            cookie: override.cookie
                        });
                        const result = this.callbackHandler({
                            cookie: override.cookie,
                            cbid: this._makecbid(),
                            get: {
                                objref,
                                property: propertyName
                            }
                        });
                        this._debug("callback returned", result);
                        return this._toSandbox(result, propInfo);
                    },
                    set: value => {
                        this._debug("virtual set", objref, propertyName, {
                            cookie: override.cookie
                        });
                        this.callbackHandler({
                            cookie: override.cookie,
                            cbid: this._makecbid(),
                            set: {
                                objref,
                                property: propertyName,
                                value: this._fromSandbox(value, propInfo)
                            }
                        });
                    }
                });
                function getPropertyDescriptor(obj, propertyName) {
                    const direct = Object.getOwnPropertyDescriptor(obj, propertyName);
                    if (direct != null) {
                        return direct;
                    }
                    const proto = Object.getPrototypeOf(obj);
                    if (proto == null && proto !== Object.prototype) {
                        return undefined;
                    }
                    return getPropertyDescriptor(proto, propertyName);
                }
            }
            _applyMethodOverride(obj, objref, typeFqn, interfaces, override) {
                if (this._tryTypeInfoForProperty(override.method, typeFqn, interfaces)) {
                    throw new Error(`Trying to override property '${override.method}' as a method`);
                }
                let methodInfo = this._tryTypeInfoForMethod(override.method, typeFqn, interfaces);
                if (!methodInfo && obj[override.method]) {
                    this._debug(`Skipping override of private method ${override.method}`);
                    return;
                }
                if (!methodInfo) {
                    methodInfo = {
                        name: override.method,
                        returns: {
                            type: spec.CANONICAL_ANY
                        },
                        parameters: [ {
                            name: "args",
                            type: spec.CANONICAL_ANY,
                            variadic: true
                        } ],
                        variadic: true
                    };
                }
                this._defineOverridenMethod(obj, objref, override, methodInfo);
            }
            _defineOverridenMethod(obj, objref, override, methodInfo) {
                const methodName = override.method;
                if (methodInfo.async) {
                    Object.defineProperty(obj, methodName, {
                        enumerable: false,
                        configurable: false,
                        writable: false,
                        value: (...methodArgs) => {
                            this._debug("invoke async method override", override);
                            const args = this._toSandboxValues(methodArgs, methodInfo.parameters);
                            return new Promise(((succeed, fail) => {
                                var _a;
                                const cbid = this._makecbid();
                                this._debug("adding callback to queue", cbid);
                                this.cbs[cbid] = {
                                    objref,
                                    override,
                                    args,
                                    expectedReturnType: (_a = methodInfo.returns) !== null && _a !== void 0 ? _a : "void",
                                    succeed,
                                    fail
                                };
                            }));
                        }
                    });
                } else {
                    Object.defineProperty(obj, methodName, {
                        enumerable: false,
                        configurable: false,
                        writable: false,
                        value: (...methodArgs) => {
                            var _a;
                            this._debug("invoke sync method override", override, "args", methodArgs);
                            const result = this.callbackHandler({
                                cookie: override.cookie,
                                cbid: this._makecbid(),
                                invoke: {
                                    objref,
                                    method: methodName,
                                    args: this._fromSandboxValues(methodArgs, methodInfo.parameters)
                                }
                            });
                            this._debug("Result", result);
                            return this._toSandbox(result, (_a = methodInfo.returns) !== null && _a !== void 0 ? _a : "void");
                        }
                    });
                }
            }
            _findInvokeTarget(objref, methodName, args) {
                const {instance, fqn, interfaces} = this.objects.findObject(objref);
                const ti = this._typeInfoForMethod(methodName, fqn, interfaces);
                this._validateMethodArguments(ti, args);
                let fn = instance.constructor.prototype[methodName];
                if (!fn) {
                    fn = instance[methodName];
                    if (!fn) {
                        throw new Error(`Cannot find ${methodName} on object`);
                    }
                }
                return {
                    ti,
                    obj: instance,
                    fn
                };
            }
            _validateMethodArguments(method, args) {
                var _a;
                const params = (_a = method === null || method === void 0 ? void 0 : method.parameters) !== null && _a !== void 0 ? _a : [];
                if (args.length > params.length && !(method && method.variadic)) {
                    throw new Error(`Too many arguments (method accepts ${params.length} parameters, got ${args.length} arguments)`);
                }
                for (let i = 0; i < params.length; ++i) {
                    const param = params[i];
                    const arg = args[i];
                    if (param.variadic) {
                        if (params.length <= i) {
                            return;
                        }
                        for (let j = i; j < params.length; j++) {
                            if (!param.optional && params[j] === undefined) {
                                throw new Error(`Unexpected 'undefined' value at index ${j - i} of variadic argument '${param.name}' of type '${spec.describeTypeReference(param.type)}'`);
                            }
                        }
                    } else if (!param.optional && arg === undefined) {
                        throw new Error(`Not enough arguments. Missing argument for the required parameter '${param.name}' of type '${spec.describeTypeReference(param.type)}'`);
                    }
                }
            }
            _assemblyFor(assemblyName) {
                const assembly = this.assemblies[assemblyName];
                if (!assembly) {
                    throw new Error(`Could not find assembly: ${assemblyName}`);
                }
                return assembly;
            }
            _findSymbol(fqn) {
                const [assemblyName, ...parts] = fqn.split(".");
                const assembly = this._assemblyFor(assemblyName);
                let curr = assembly.closure;
                while (parts.length > 0) {
                    const name = parts.shift();
                    if (!name) {
                        break;
                    }
                    curr = curr[name];
                }
                if (!curr) {
                    throw new Error(`Could not find symbol ${fqn}`);
                }
                return curr;
            }
            _typeInfoForFqn(fqn) {
                var _a;
                const components = fqn.split(".");
                const moduleName = components[0];
                const assembly = this.assemblies[moduleName];
                if (!assembly) {
                    throw new Error(`Module '${moduleName}' not found`);
                }
                const types = (_a = assembly.metadata.types) !== null && _a !== void 0 ? _a : {};
                const fqnInfo = types[fqn];
                if (!fqnInfo) {
                    throw new Error(`Type '${fqn}' not found`);
                }
                return fqnInfo;
            }
            _typeInfoForMethod(methodName, fqn, interfaces) {
                const ti = this._tryTypeInfoForMethod(methodName, fqn, interfaces);
                if (!ti) {
                    const addendum = interfaces && interfaces.length > 0 ? ` or interface(s) ${interfaces.join(", ")}` : "";
                    throw new Error(`Class ${fqn}${addendum} doesn't have a method '${methodName}'`);
                }
                return ti;
            }
            _tryTypeInfoForMethod(methodName, classFqn, interfaces = []) {
                var _a, _b;
                for (const fqn of [ classFqn, ...interfaces ]) {
                    if (fqn === "Object") {
                        continue;
                    }
                    const typeinfo = this._typeInfoForFqn(fqn);
                    const methods = (_a = typeinfo.methods) !== null && _a !== void 0 ? _a : [];
                    for (const m of methods) {
                        if (m.name === methodName) {
                            return m;
                        }
                    }
                    const bases = [ typeinfo.base, ...(_b = typeinfo.interfaces) !== null && _b !== void 0 ? _b : [] ];
                    for (const base of bases) {
                        if (!base) {
                            continue;
                        }
                        const found = this._tryTypeInfoForMethod(methodName, base);
                        if (found) {
                            return found;
                        }
                    }
                }
                return undefined;
            }
            _tryTypeInfoForProperty(property, classFqn, interfaces = []) {
                var _a;
                for (const fqn of [ classFqn, ...interfaces ]) {
                    if (fqn === wire.EMPTY_OBJECT_FQN) {
                        continue;
                    }
                    const typeInfo = this._typeInfoForFqn(fqn);
                    let properties;
                    let bases;
                    if (spec.isClassType(typeInfo)) {
                        const classTypeInfo = typeInfo;
                        properties = classTypeInfo.properties;
                        bases = classTypeInfo.base ? [ classTypeInfo.base ] : [];
                    } else if (spec.isInterfaceType(typeInfo)) {
                        const interfaceTypeInfo = typeInfo;
                        properties = interfaceTypeInfo.properties;
                        bases = (_a = interfaceTypeInfo.interfaces) !== null && _a !== void 0 ? _a : [];
                    } else {
                        throw new Error(`Type of kind ${typeInfo.kind} does not have properties`);
                    }
                    for (const p of properties !== null && properties !== void 0 ? properties : []) {
                        if (p.name === property) {
                            return p;
                        }
                    }
                    for (const baseFqn of bases) {
                        const ret = this._tryTypeInfoForProperty(property, baseFqn);
                        if (ret) {
                            return ret;
                        }
                    }
                }
                return undefined;
            }
            _typeInfoForProperty(property, fqn, interfaces) {
                const typeInfo = this._tryTypeInfoForProperty(property, fqn, interfaces);
                if (!typeInfo) {
                    const addendum = interfaces && interfaces.length > 0 ? ` or interface(s) ${interfaces.join(", ")}` : "";
                    throw new Error(`Type ${fqn}${addendum} doesn't have a property '${property}'`);
                }
                return typeInfo;
            }
            _toSandbox(v, expectedType) {
                const serTypes = wire.serializationType(expectedType, this._typeInfoForFqn.bind(this));
                this._debug("toSandbox", v, JSON.stringify(serTypes));
                const host = {
                    objects: this.objects,
                    debug: this._debug.bind(this),
                    findSymbol: this._findSymbol.bind(this),
                    lookupType: this._typeInfoForFqn.bind(this),
                    recurse: this._toSandbox.bind(this)
                };
                const errors = new Array;
                for (const {serializationClass, typeRef} of serTypes) {
                    try {
                        return wire.SERIALIZERS[serializationClass].deserialize(v, typeRef, host);
                    } catch (e) {
                        if (serTypes.length === 1) {
                            throw e;
                        }
                        errors.push(e.message);
                    }
                }
                throw new Error(`Value did not match any type in union: ${errors.join(", ")}`);
            }
            _fromSandbox(v, targetType) {
                const serTypes = wire.serializationType(targetType, this._typeInfoForFqn.bind(this));
                this._debug("fromSandbox", v, JSON.stringify(serTypes));
                const host = {
                    objects: this.objects,
                    debug: this._debug.bind(this),
                    findSymbol: this._findSymbol.bind(this),
                    lookupType: this._typeInfoForFqn.bind(this),
                    recurse: this._fromSandbox.bind(this)
                };
                const errors = new Array;
                for (const {serializationClass, typeRef} of serTypes) {
                    try {
                        return wire.SERIALIZERS[serializationClass].serialize(v, typeRef, host);
                    } catch (e) {
                        if (serTypes.length === 1) {
                            throw e;
                        }
                        errors.push(e.message);
                    }
                }
                throw new Error(`Value did not match any type in union: ${errors.join(", ")}`);
            }
            _toSandboxValues(xs, parameters) {
                return this._boxUnboxParameters(xs, parameters, this._toSandbox.bind(this));
            }
            _fromSandboxValues(xs, parameters) {
                return this._boxUnboxParameters(xs, parameters, this._fromSandbox.bind(this));
            }
            _boxUnboxParameters(xs, parameters, boxUnbox) {
                parameters = [ ...parameters !== null && parameters !== void 0 ? parameters : [] ];
                const variadic = parameters.length > 0 && !!parameters[parameters.length - 1].variadic;
                while (variadic && parameters.length < xs.length) {
                    parameters.push(parameters[parameters.length - 1]);
                }
                if (xs.length > parameters.length) {
                    throw new Error(`Argument list (${JSON.stringify(xs)}) not same size as expected argument list (length ${parameters.length})`);
                }
                return xs.map(((x, i) => boxUnbox(x, parameters[i])));
            }
            _debug(...args) {
                if (this.traceEnabled) {
                    console.error("[@jsii/kernel]", ...args);
                }
            }
            _ensureSync(desc, fn) {
                this.syncInProgress = desc;
                try {
                    return fn();
                } finally {
                    delete this.syncInProgress;
                }
            }
            _findPropertyTarget(obj, property) {
                const superProp = this._getSuperPropertyName(property);
                if (superProp in obj) {
                    return superProp;
                }
                return property;
            }
            _makecbid() {
                return `jsii::callback::${this.nextid++}`;
            }
            _makeprid() {
                return `jsii::promise::${this.nextid++}`;
            }
            _wrapSandboxCode(fn) {
                return fn();
            }
            _execute(code, filename) {
                const script = new vm.Script(code, {
                    filename
                });
                return script.runInContext(this.sandbox, {
                    displayErrors: true
                });
            }
        }
        exports.Kernel = Kernel;
        class Assembly {
            constructor(metadata, closure) {
                this.metadata = metadata;
                this.closure = closure;
            }
        }
    },
    2309: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.ObjectTable = exports.tagJsiiConstructor = exports.objectReference = exports.jsiiTypeFqn = void 0;
        const spec = __webpack_require__(1804);
        const api = __webpack_require__(2816);
        const serialization_1 = __webpack_require__(3837);
        const OBJID_SYMBOL = Symbol.for("$__jsii__objid__$");
        const IFACES_SYMBOL = Symbol.for("$__jsii__interfaces__$");
        const JSII_SYMBOL = Symbol.for("__jsii__");
        function jsiiTypeFqn(obj) {
            var _a;
            return (_a = obj.constructor[JSII_SYMBOL]) === null || _a === void 0 ? void 0 : _a.fqn;
        }
        exports.jsiiTypeFqn = jsiiTypeFqn;
        function objectReference(obj) {
            if (obj[OBJID_SYMBOL]) {
                return {
                    [api.TOKEN_REF]: obj[OBJID_SYMBOL],
                    [api.TOKEN_INTERFACES]: obj[IFACES_SYMBOL]
                };
            }
            return undefined;
        }
        exports.objectReference = objectReference;
        function tagObject(obj, objid, interfaces) {
            const managed = obj;
            managed[OBJID_SYMBOL] = objid;
            managed[IFACES_SYMBOL] = interfaces;
        }
        function tagJsiiConstructor(constructor, fqn) {
            Object.defineProperty(constructor, JSII_SYMBOL, {
                configurable: false,
                enumerable: false,
                writable: false,
                value: {
                    fqn
                }
            });
        }
        exports.tagJsiiConstructor = tagJsiiConstructor;
        class ObjectTable {
            constructor(resolveType) {
                this.resolveType = resolveType;
                this.objects = {};
                this.nextid = 1e4;
            }
            registerObject(obj, fqn, interfaces) {
                var _a;
                if (fqn === undefined) {
                    throw new Error("FQN cannot be undefined");
                }
                const existingRef = objectReference(obj);
                if (existingRef) {
                    if (interfaces) {
                        const allIfaces = new Set(interfaces);
                        for (const iface of (_a = existingRef[api.TOKEN_INTERFACES]) !== null && _a !== void 0 ? _a : []) {
                            allIfaces.add(iface);
                        }
                        this.objects[existingRef[api.TOKEN_REF]].interfaces = obj[IFACES_SYMBOL] = existingRef[api.TOKEN_INTERFACES] = interfaces = this.removeRedundant(Array.from(allIfaces), fqn);
                    }
                    return existingRef;
                }
                interfaces = this.removeRedundant(interfaces, fqn);
                const objid = this.makeId(fqn);
                this.objects[objid] = {
                    instance: obj,
                    fqn,
                    interfaces
                };
                tagObject(obj, objid, interfaces);
                return {
                    [api.TOKEN_REF]: objid,
                    [api.TOKEN_INTERFACES]: interfaces
                };
            }
            findObject(objref) {
                if (typeof objref !== "object" || !(api.TOKEN_REF in objref)) {
                    throw new Error(`Malformed object reference: ${JSON.stringify(objref)}`);
                }
                const objid = objref[api.TOKEN_REF];
                const obj = this.objects[objid];
                if (!obj) {
                    throw new Error(`Object ${objid} not found`);
                }
                return obj;
            }
            deleteObject(objref) {
                this.findObject(objref);
                delete this.objects[objref[api.TOKEN_REF]];
            }
            get count() {
                return Object.keys(this.objects).length;
            }
            makeId(fqn) {
                return `${fqn}@${this.nextid++}`;
            }
            removeRedundant(interfaces, fqn) {
                if (!interfaces || interfaces.length === 0) {
                    return undefined;
                }
                const result = new Set(interfaces);
                const builtIn = new InterfaceCollection(this.resolveType);
                if (fqn !== serialization_1.EMPTY_OBJECT_FQN) {
                    builtIn.addFromClass(fqn);
                }
                interfaces.forEach(builtIn.addFromInterface.bind(builtIn));
                for (const iface of builtIn) {
                    result.delete(iface);
                }
                return result.size > 0 ? Array.from(result).sort() : undefined;
            }
        }
        exports.ObjectTable = ObjectTable;
        class InterfaceCollection {
            constructor(resolveType) {
                this.resolveType = resolveType;
                this.interfaces = new Set;
            }
            addFromClass(fqn) {
                const ti = this.resolveType(fqn);
                if (!spec.isClassType(ti)) {
                    throw new Error(`Expected a class, but received ${spec.describeTypeReference(ti)}`);
                }
                if (ti.base) {
                    this.addFromClass(ti.base);
                }
                if (ti.interfaces) {
                    for (const iface of ti.interfaces) {
                        if (this.interfaces.has(iface)) {
                            continue;
                        }
                        this.interfaces.add(iface);
                        this.addFromInterface(iface);
                    }
                }
            }
            addFromInterface(fqn) {
                const ti = this.resolveType(fqn);
                if (!spec.isInterfaceType(ti)) {
                    throw new Error(`Expected an interface, but received ${spec.describeTypeReference(ti)}`);
                }
                if (!ti.interfaces) {
                    return;
                }
                for (const iface of ti.interfaces) {
                    if (this.interfaces.has(iface)) {
                        continue;
                    }
                    this.interfaces.add(iface);
                    this.addFromInterface(iface);
                }
            }
            [Symbol.iterator]() {
                return this.interfaces[Symbol.iterator]();
            }
        }
    },
    6703: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.removeSync = void 0;
        const fs = __webpack_require__(9728);
        const process = __webpack_require__(1765);
        const removeSyncPaths = new Array;
        function removeSync(path) {
            registerIfNeeded();
            removeSyncPaths.push(path);
        }
        exports.removeSync = removeSync;
        let registered = false;
        function registerIfNeeded() {
            if (registered) {
                return;
            }
            process.once("exit", onExitHandler);
            registered = true;
            function onExitHandler() {
                if (removeSyncPaths.length > 0) {
                    for (const path of removeSyncPaths) {
                        fs.removeSync(path);
                    }
                }
            }
        }
    },
    3837: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.serializationType = exports.SERIALIZERS = exports.SYMBOL_WIRE_TYPE = exports.EMPTY_OBJECT_FQN = void 0;
        const spec = __webpack_require__(1804);
        const api_1 = __webpack_require__(2816);
        const objects_1 = __webpack_require__(2309);
        const _1 = __webpack_require__(8944);
        exports.EMPTY_OBJECT_FQN = "Object";
        exports.SYMBOL_WIRE_TYPE = Symbol.for("$jsii$wireType$");
        exports.SERIALIZERS = {
            ["Void"]: {
                serialize(value, _type, host) {
                    if (value != null) {
                        host.debug("Expected void, got", value);
                    }
                    return undefined;
                },
                deserialize(value, _type, host) {
                    if (value != null) {
                        host.debug("Expected void, got", value);
                    }
                    return undefined;
                }
            },
            ["Date"]: {
                serialize(value, optionalValue) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    if (!isDate(value)) {
                        throw new Error(`Expected Date, got ${JSON.stringify(value)}`);
                    }
                    return serializeDate(value);
                },
                deserialize(value, optionalValue) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (!api_1.isWireDate(value)) {
                        throw new Error(`Expected Date, got ${JSON.stringify(value)}`);
                    }
                    return deserializeDate(value);
                }
            },
            ["Scalar"]: {
                serialize(value, optionalValue) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    const primitiveType = optionalValue.type;
                    if (!isScalar(value)) {
                        throw new Error(`Expected ${spec.describeTypeReference(optionalValue.type)}, got ${JSON.stringify(value)}`);
                    }
                    if (typeof value !== primitiveType.primitive) {
                        throw new Error(`Expected a ${spec.describeTypeReference(optionalValue.type)}, got ${JSON.stringify(value)} (${typeof value})`);
                    }
                    return value;
                },
                deserialize(value, optionalValue) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    const primitiveType = optionalValue.type;
                    if (!isScalar(value)) {
                        throw new Error(`Expected a ${spec.describeTypeReference(optionalValue.type)}, got ${JSON.stringify(value)}`);
                    }
                    if (typeof value !== primitiveType.primitive) {
                        throw new Error(`Expected a ${spec.describeTypeReference(optionalValue.type)}, got ${JSON.stringify(value)} (${typeof value})`);
                    }
                    return value;
                }
            },
            ["Json"]: {
                serialize(value) {
                    return value;
                },
                deserialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (api_1.isWireMap(value)) {
                        return exports.SERIALIZERS["Map"].deserialize(value, {
                            optional: false,
                            type: {
                                collection: {
                                    kind: spec.CollectionKind.Map,
                                    elementtype: {
                                        primitive: spec.PrimitiveType.Json
                                    }
                                }
                            }
                        }, host);
                    }
                    if (typeof value !== "object") {
                        return value;
                    }
                    if (Array.isArray(value)) {
                        return value.map(mapJsonValue);
                    }
                    return mapValues(value, mapJsonValue);
                    function mapJsonValue(toMap) {
                        if (toMap == null) {
                            return toMap;
                        }
                        return host.recurse(toMap, {
                            type: {
                                primitive: spec.PrimitiveType.Json
                            }
                        });
                    }
                }
            },
            ["Enum"]: {
                serialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    if (typeof value !== "string" && typeof value !== "number") {
                        throw new Error(`Expected enum value, got ${JSON.stringify(value)}`);
                    }
                    host.debug("Serializing enum");
                    const enumType = optionalValue.type;
                    const enumMap = host.findSymbol(enumType.fqn);
                    const enumEntry = Object.entries(enumMap).find((([, v]) => v === value));
                    if (!enumEntry) {
                        throw new Error(`No entry in ${enumType.fqn} has value ${value}`);
                    }
                    return {
                        [api_1.TOKEN_ENUM]: `${enumType.fqn}/${enumEntry[0]}`
                    };
                },
                deserialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (!api_1.isWireEnum(value)) {
                        throw new Error(`Expected enum value, got ${JSON.stringify(value)}`);
                    }
                    return deserializeEnum(value, host.findSymbol);
                }
            },
            ["Array"]: {
                serialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    if (!Array.isArray(value)) {
                        throw new Error(`Expected array type, got ${JSON.stringify(value)}`);
                    }
                    const arrayType = optionalValue.type;
                    return value.map((x => host.recurse(x, {
                        type: arrayType.collection.elementtype
                    })));
                },
                deserialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    if (!Array.isArray(value)) {
                        throw new Error(`Expected array type, got ${JSON.stringify(value)}`);
                    }
                    const arrayType = optionalValue.type;
                    return value.map((x => host.recurse(x, {
                        type: arrayType.collection.elementtype
                    })));
                }
            },
            ["Map"]: {
                serialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    const mapType = optionalValue.type;
                    return {
                        [api_1.TOKEN_MAP]: mapValues(value, (v => host.recurse(v, {
                            type: mapType.collection.elementtype
                        })))
                    };
                },
                deserialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    const mapType = optionalValue.type;
                    if (!api_1.isWireMap(value)) {
                        return mapValues(value, (v => host.recurse(v, {
                            type: mapType.collection.elementtype
                        })));
                    }
                    const result = mapValues(value[api_1.TOKEN_MAP], (v => host.recurse(v, {
                        type: mapType.collection.elementtype
                    })));
                    Object.defineProperty(result, exports.SYMBOL_WIRE_TYPE, {
                        configurable: false,
                        enumerable: false,
                        value: api_1.TOKEN_MAP,
                        writable: false
                    });
                    return result;
                }
            },
            ["Struct"]: {
                serialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    if (typeof value !== "object" || value == null) {
                        throw new Error(`Expected object, got ${JSON.stringify(value)}`);
                    }
                    host.debug("Returning value type by reference");
                    return host.objects.registerObject(value, "Object", [ optionalValue.type.fqn ]);
                },
                deserialize(value, optionalValue, host) {
                    if (typeof value === "object" && Object.keys(value !== null && value !== void 0 ? value : {}).length === 0) {
                        value = undefined;
                    }
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    if (typeof value !== "object" || value == null) {
                        throw new Error(`Expected object reference, got ${JSON.stringify(value)}`);
                    }
                    const namedType = host.lookupType(optionalValue.type.fqn);
                    const props = propertiesOf(namedType, host.lookupType);
                    if (Array.isArray(value)) {
                        throw new Error(`Got an array where a ${namedType.fqn} was expected. Did you mean to pass a variable number of arguments?`);
                    }
                    if (api_1.isObjRef(value)) {
                        host.debug("Expected value type but got reference type, accepting for now (awslabs/jsii#400)");
                        return validateRequiredProps(host.objects.findObject(value).instance, namedType.fqn, props);
                    }
                    if (_1.api.isWireStruct(value)) {
                        const {fqn, data} = value[_1.api.TOKEN_STRUCT];
                        if (!isAssignable(fqn, namedType, host.lookupType)) {
                            throw new Error(`Wire struct type '${fqn}' does not match expected '${namedType.fqn}'`);
                        }
                        value = data;
                    }
                    if (_1.api.isWireMap(value)) {
                        value = value[_1.api.TOKEN_MAP];
                    }
                    value = validateRequiredProps(value, namedType.fqn, props);
                    return mapValues(value, ((v, key) => {
                        if (!props[key]) {
                            return undefined;
                        }
                        return host.recurse(v, props[key]);
                    }));
                }
            },
            ["RefType"]: {
                serialize(value, optionalValue, host) {
                    var _a;
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    if (typeof value !== "object" || value == null) {
                        throw new Error(`Expected object reference, got ${JSON.stringify(value)}`);
                    }
                    const expectedType = host.lookupType(optionalValue.type.fqn);
                    const interfaces = spec.isInterfaceType(expectedType) ? [ expectedType.fqn ] : undefined;
                    const jsiiType = (_a = objects_1.jsiiTypeFqn(value)) !== null && _a !== void 0 ? _a : spec.isClassType(expectedType) ? expectedType.fqn : "Object";
                    return host.objects.registerObject(value, jsiiType, interfaces);
                },
                deserialize(value, optionalValue, host) {
                    if (nullAndOk(value, optionalValue)) {
                        return undefined;
                    }
                    if (optionalValue === "void") {
                        throw new Error("Encountered unexpected `void` type");
                    }
                    if (!api_1.isObjRef(value)) {
                        throw new Error(`Expected object reference, got ${JSON.stringify(value)}`);
                    }
                    const {instance, fqn} = host.objects.findObject(value);
                    const namedTypeRef = optionalValue.type;
                    if (namedTypeRef.fqn !== exports.EMPTY_OBJECT_FQN) {
                        const namedType = host.lookupType(namedTypeRef.fqn);
                        const declaredType = optionalValue.type;
                        if (spec.isClassType(namedType) && !isAssignable(fqn, declaredType, host.lookupType)) {
                            throw new Error(`Object of type ${fqn} is not convertible to ${declaredType.fqn}`);
                        }
                    }
                    return instance;
                }
            },
            ["Any"]: {
                serialize(value, _type, host) {
                    var _a;
                    if (value == null) {
                        return undefined;
                    }
                    if (isDate(value)) {
                        return serializeDate(value);
                    }
                    if (isScalar(value)) {
                        return value;
                    }
                    if (Array.isArray(value)) {
                        return value.map((e => host.recurse(e, {
                            type: spec.CANONICAL_ANY
                        })));
                    }
                    if (typeof value === "function") {
                        throw new Error("JSII Kernel is unable to serialize `function`. An instance with methods might have been returned by an `any` method?");
                    }
                    if (typeof value !== "object" || value == null) {
                        throw new Error(`JSII kernel assumption violated, ${JSON.stringify(value)} is not an object`);
                    }
                    if (exports.SYMBOL_WIRE_TYPE in value && value[exports.SYMBOL_WIRE_TYPE] === api_1.TOKEN_MAP) {
                        return exports.SERIALIZERS["Map"].serialize(value, {
                            type: {
                                collection: {
                                    kind: spec.CollectionKind.Map,
                                    elementtype: spec.CANONICAL_ANY
                                }
                            }
                        }, host);
                    }
                    if (value instanceof Set || value instanceof Map) {
                        throw new Error("Can't return objects of type Set or Map");
                    }
                    const prevRef = objects_1.objectReference(value);
                    if (prevRef) {
                        return prevRef;
                    }
                    const jsiiType = (_a = objects_1.jsiiTypeFqn(value)) !== null && _a !== void 0 ? _a : isByReferenceOnly(value) ? exports.EMPTY_OBJECT_FQN : undefined;
                    if (jsiiType) {
                        return host.objects.registerObject(value, jsiiType);
                    }
                    return mapValues(value, (v => host.recurse(v, {
                        type: spec.CANONICAL_ANY
                    })));
                },
                deserialize(value, _type, host) {
                    if (value == null) {
                        return undefined;
                    }
                    if (api_1.isWireDate(value)) {
                        host.debug("ANY is a Date");
                        return deserializeDate(value);
                    }
                    if (isScalar(value)) {
                        host.debug("ANY is a Scalar");
                        return value;
                    }
                    if (Array.isArray(value)) {
                        host.debug("ANY is an Array");
                        return value.map((e => host.recurse(e, {
                            type: spec.CANONICAL_ANY
                        })));
                    }
                    if (api_1.isWireEnum(value)) {
                        host.debug("ANY is an Enum");
                        return deserializeEnum(value, host.findSymbol);
                    }
                    if (api_1.isWireMap(value)) {
                        host.debug("ANY is a Map");
                        const mapOfAny = {
                            collection: {
                                kind: spec.CollectionKind.Map,
                                elementtype: spec.CANONICAL_ANY
                            }
                        };
                        return exports.SERIALIZERS["Map"].deserialize(value, {
                            type: mapOfAny
                        }, host);
                    }
                    if (api_1.isObjRef(value)) {
                        host.debug("ANY is a Ref");
                        return host.objects.findObject(value).instance;
                    }
                    if (api_1.isWireStruct(value)) {
                        const {fqn, data} = value[api_1.TOKEN_STRUCT];
                        host.debug(`ANY is a struct of type ${fqn}`);
                        return exports.SERIALIZERS["Struct"].deserialize(data, {
                            type: {
                                fqn
                            }
                        }, host);
                    }
                    host.debug("ANY is a Map");
                    return mapValues(value, (v => host.recurse(v, {
                        type: spec.CANONICAL_ANY
                    })));
                }
            }
        };
        function serializeDate(value) {
            return {
                [api_1.TOKEN_DATE]: value.toISOString()
            };
        }
        function deserializeDate(value) {
            return new Date(value[api_1.TOKEN_DATE]);
        }
        function deserializeEnum(value, lookup) {
            const enumLocator = value[api_1.TOKEN_ENUM];
            const sep = enumLocator.lastIndexOf("/");
            if (sep === -1) {
                throw new Error(`Malformed enum value: ${JSON.stringify(value)}`);
            }
            const typeName = enumLocator.substr(0, sep);
            const valueName = enumLocator.substr(sep + 1);
            const enumValue = lookup(typeName)[valueName];
            if (enumValue === undefined) {
                throw new Error(`No enum member named ${valueName} in ${typeName}`);
            }
            return enumValue;
        }
        function serializationType(typeRef, lookup) {
            if (typeRef == null) {
                throw new Error("Kernel error: expected type information, got 'undefined'");
            }
            if (typeRef === "void") {
                return [ {
                    serializationClass: "Void",
                    typeRef
                } ];
            }
            if (spec.isPrimitiveTypeReference(typeRef.type)) {
                switch (typeRef.type.primitive) {
                  case spec.PrimitiveType.Any:
                    return [ {
                        serializationClass: "Any",
                        typeRef
                    } ];

                  case spec.PrimitiveType.Date:
                    return [ {
                        serializationClass: "Date",
                        typeRef
                    } ];

                  case spec.PrimitiveType.Json:
                    return [ {
                        serializationClass: "Json",
                        typeRef
                    } ];

                  case spec.PrimitiveType.Boolean:
                  case spec.PrimitiveType.Number:
                  case spec.PrimitiveType.String:
                    return [ {
                        serializationClass: "Scalar",
                        typeRef
                    } ];
                }
                throw new Error("Unknown primitive type");
            }
            if (spec.isCollectionTypeReference(typeRef.type)) {
                return [ {
                    serializationClass: typeRef.type.collection.kind === spec.CollectionKind.Array ? "Array" : "Map",
                    typeRef
                } ];
            }
            if (spec.isUnionTypeReference(typeRef.type)) {
                const compoundTypes = flatMap(typeRef.type.union.types, (t => serializationType({
                    type: t
                }, lookup)));
                for (const t of compoundTypes) {
                    if (t.typeRef !== "void") {
                        t.typeRef.optional = typeRef.optional;
                    }
                }
                return compoundTypes.sort(((l, r) => compareSerializationClasses(l.serializationClass, r.serializationClass)));
            }
            const type = lookup(typeRef.type.fqn);
            if (spec.isEnumType(type)) {
                return [ {
                    serializationClass: "Enum",
                    typeRef
                } ];
            }
            if (spec.isInterfaceType(type) && type.datatype) {
                return [ {
                    serializationClass: "Struct",
                    typeRef
                } ];
            }
            return [ {
                serializationClass: "RefType",
                typeRef
            } ];
        }
        exports.serializationType = serializationType;
        function nullAndOk(x, type) {
            if (x != null) {
                return false;
            }
            if (type !== "void" && !type.optional) {
                throw new Error(`Got 'undefined' for non-optional instance of ${JSON.stringify(type)}`);
            }
            return true;
        }
        function isDate(x) {
            return typeof x === "object" && Object.prototype.toString.call(x) === "[object Date]";
        }
        function isScalar(x) {
            return typeof x === "string" || typeof x === "number" || typeof x === "boolean";
        }
        function flatMap(xs, fn) {
            const ret = new Array;
            for (const x of xs) {
                ret.push(...fn(x));
            }
            return ret;
        }
        function mapValues(value, fn) {
            if (typeof value !== "object" || value == null) {
                throw new Error(`Expected object type, got ${JSON.stringify(value)}`);
            }
            const out = {};
            for (const [k, v] of Object.entries(value)) {
                const wireValue = fn(v, k);
                if (wireValue === undefined) {
                    continue;
                }
                out[k] = wireValue;
            }
            return out;
        }
        function propertiesOf(t, lookup) {
            var _a;
            if (!spec.isClassOrInterfaceType(t)) {
                return {};
            }
            let ret = {};
            if (t.interfaces) {
                for (const iface of t.interfaces) {
                    ret = {
                        ...ret,
                        ...propertiesOf(lookup(iface), lookup)
                    };
                }
            }
            if (spec.isClassType(t) && t.base) {
                ret = {
                    ...ret,
                    ...propertiesOf(lookup(t.base), lookup)
                };
            }
            for (const prop of (_a = t.properties) !== null && _a !== void 0 ? _a : []) {
                ret[prop.name] = prop;
            }
            return ret;
        }
        function isAssignable(actualTypeFqn, requiredType, lookup) {
            if (actualTypeFqn === exports.EMPTY_OBJECT_FQN) {
                return true;
            }
            if (requiredType.fqn === actualTypeFqn) {
                return true;
            }
            const actualType = lookup(actualTypeFqn);
            if (spec.isClassType(actualType)) {
                if (actualType.base && isAssignable(actualType.base, requiredType, lookup)) {
                    return true;
                }
            }
            if (spec.isClassOrInterfaceType(actualType) && actualType.interfaces) {
                return actualType.interfaces.find((iface => isAssignable(iface, requiredType, lookup))) != null;
            }
            return false;
        }
        function validateRequiredProps(actualProps, typeName, specProps) {
            const missingRequiredProps = Object.keys(specProps).filter((name => !specProps[name].optional)).filter((name => !(name in actualProps)));
            if (missingRequiredProps.length > 0) {
                throw new Error(`Missing required properties for ${typeName}: ${missingRequiredProps.join(", ")}`);
            }
            return actualProps;
        }
        function compareSerializationClasses(l, r) {
            const order = [ "Void", "Date", "Scalar", "Json", "Enum", "Array", "Map", "Struct", "RefType", "Any" ];
            return order.indexOf(l) - order.indexOf(r);
        }
        function isByReferenceOnly(obj) {
            if (Array.isArray(obj)) {
                return false;
            }
            let curr = obj;
            do {
                for (const prop of Object.getOwnPropertyNames(curr)) {
                    const descr = Object.getOwnPropertyDescriptor(curr, prop);
                    if ((descr === null || descr === void 0 ? void 0 : descr.get) != null || (descr === null || descr === void 0 ? void 0 : descr.set) != null || typeof (descr === null || descr === void 0 ? void 0 : descr.value) === "function") {
                        return true;
                    }
                }
            } while (Object.getPrototypeOf(curr = Object.getPrototypeOf(curr)) != null);
            return false;
        }
    },
    7905: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.KernelHost = void 0;
        const kernel_1 = __webpack_require__(8944);
        const events_1 = __webpack_require__(8614);
        class KernelHost {
            constructor(inout, opts = {}) {
                this.inout = inout;
                this.opts = opts;
                this.kernel = new kernel_1.Kernel(this.callbackHandler.bind(this));
                this.eventEmitter = new events_1.EventEmitter;
                this.kernel.traceEnabled = opts.debug ? true : false;
            }
            run() {
                var _a;
                const req = this.inout.read();
                if (!req || "exit" in req) {
                    this.eventEmitter.emit("exit", (_a = req === null || req === void 0 ? void 0 : req.exit) !== null && _a !== void 0 ? _a : 0);
                    return;
                }
                this.processRequest(req, (() => {
                    setImmediate((() => this.run()));
                }));
            }
            once(event, listener) {
                this.eventEmitter.once(event, listener);
            }
            callbackHandler(callback) {
                this.inout.write({
                    callback
                });
                return completeCallback.call(this);
                function completeCallback() {
                    const req = this.inout.read();
                    if (!req || "exit" in req) {
                        throw new Error("Interrupted before callback returned");
                    }
                    const completeReq = req;
                    if ("complete" in completeReq && completeReq.complete.cbid === callback.cbid) {
                        if (completeReq.complete.err) {
                            throw new Error(completeReq.complete.err);
                        }
                        return completeReq.complete.result;
                    }
                    return this.processRequest(req, completeCallback.bind(this), true);
                }
            }
            processRequest(req, next, sync = false) {
                if ("callback" in req) {
                    throw new Error("Unexpected `callback` result. This request should have been processed by a callback handler");
                }
                if (!("api" in req)) {
                    throw new Error('Malformed request, "api" field is required');
                }
                const apiReq = req;
                const fn = this.findApi(apiReq.api);
                try {
                    const ret = fn.call(this.kernel, req);
                    if (apiReq.api === "begin" || apiReq.api === "complete") {
                        checkIfAsyncIsAllowed();
                        this.debug("processing pending promises before responding");
                        setImmediate((() => {
                            this.writeOkay(ret);
                            next();
                        }));
                        return undefined;
                    }
                    if (this.isPromise(ret)) {
                        checkIfAsyncIsAllowed();
                        this.debug("waiting for promise to be fulfilled");
                        const promise = ret;
                        promise.then((val => {
                            this.debug("promise succeeded:", val);
                            this.writeOkay(val);
                            next();
                        })).catch((e => {
                            this.debug("promise failed:", e);
                            this.writeError(e);
                            next();
                        }));
                        return undefined;
                    }
                    this.writeOkay(ret);
                } catch (e) {
                    this.writeError(e);
                }
                return next();
                function checkIfAsyncIsAllowed() {
                    if (sync) {
                        throw new Error("Cannot handle async operations while waiting for a sync callback to return");
                    }
                }
            }
            writeOkay(result) {
                const res = {
                    ok: result
                };
                this.inout.write(res);
            }
            writeError(error) {
                const res = {
                    error: error.message,
                    stack: undefined
                };
                if (!this.opts.noStack) {
                    res.stack = error.stack;
                }
                this.inout.write(res);
            }
            isPromise(v) {
                return typeof (v === null || v === void 0 ? void 0 : v.then) === "function";
            }
            findApi(apiName) {
                const fn = this.kernel[apiName];
                if (typeof fn !== "function") {
                    throw new Error(`Invalid kernel api call: ${apiName}`);
                }
                return fn;
            }
            debug(...args) {
                if (!this.opts.debug) {
                    return;
                }
                console.error(...args);
            }
        }
        exports.KernelHost = KernelHost;
    },
    6156: (__unused_webpack_module, exports) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.InputOutput = void 0;
        class InputOutput {
            constructor(stdio) {
                this.stdio = stdio;
                this.debug = false;
            }
            write(obj) {
                const output = JSON.stringify(obj);
                this.stdio.writeLine(output);
                if (this.debug) {
                    this.stdio.writeErrorLine(`< ${output}`);
                }
            }
            read() {
                let reqLine = this.stdio.readLine();
                if (!reqLine) {
                    return undefined;
                }
                if (reqLine.startsWith("< ")) {
                    return this.read();
                }
                if (reqLine.startsWith("> ")) {
                    reqLine = reqLine.substr(2);
                }
                const input = JSON.parse(reqLine);
                if (this.debug) {
                    this.stdio.writeErrorLine(`> ${JSON.stringify(input)}`);
                }
                return input;
            }
        }
        exports.InputOutput = InputOutput;
    },
    1416: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.SyncStdio = void 0;
        const fs = __webpack_require__(5747);
        const INPUT_BUFFER_SIZE = 1048576;
        class SyncStdio {
            constructor({errorFD, readFD, writeFD}) {
                this.bufferedData = Buffer.alloc(0);
                this.readBuffer = Buffer.alloc(INPUT_BUFFER_SIZE);
                this.stderr = errorFD;
                this.stdin = readFD;
                this.stdout = writeFD;
            }
            writeErrorLine(line) {
                this.writeBuffer(Buffer.from(`${line}\n`), this.stderr);
            }
            writeLine(line) {
                this.writeBuffer(Buffer.from(`${line}\n`), this.stdout);
            }
            readLine() {
                while (!this.bufferedData.includes("\n", 0, "utf-8")) {
                    const read = fs.readSync(this.stdin, this.readBuffer, 0, this.readBuffer.length, null);
                    if (read === 0) {
                        return undefined;
                    }
                    const newData = this.readBuffer.slice(0, read);
                    this.bufferedData = Buffer.concat([ this.bufferedData, newData ]);
                }
                const newLinePos = this.bufferedData.indexOf("\n", 0, "utf-8");
                const next = this.bufferedData.slice(0, newLinePos).toString("utf-8");
                this.bufferedData = this.bufferedData.slice(newLinePos + 1);
                return next;
            }
            writeBuffer(buffer, fd) {
                let offset = 0;
                while (offset < buffer.length) {
                    try {
                        offset += fs.writeSync(fd, buffer, offset);
                    } catch (e) {
                        if (e.code !== "EAGAIN") {
                            throw e;
                        }
                    }
                }
            }
        }
        exports.SyncStdio = SyncStdio;
    },
    2752: (__unused_webpack_module, exports) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.isDeprecated = exports.describeTypeReference = exports.isClassOrInterfaceType = exports.isEnumType = exports.isInterfaceType = exports.isClassType = exports.TypeKind = exports.isMethod = exports.isUnionTypeReference = exports.isCollectionTypeReference = exports.isPrimitiveTypeReference = exports.isNamedTypeReference = exports.CANONICAL_ANY = exports.PrimitiveType = exports.CollectionKind = exports.Stability = exports.SchemaVersion = exports.SPEC_FILE_NAME = void 0;
        exports.SPEC_FILE_NAME = ".jsii";
        var SchemaVersion;
        (function(SchemaVersion) {
            SchemaVersion["LATEST"] = "jsii/0.10.0";
        })(SchemaVersion = exports.SchemaVersion || (exports.SchemaVersion = {}));
        var Stability;
        (function(Stability) {
            Stability["Deprecated"] = "deprecated";
            Stability["Experimental"] = "experimental";
            Stability["Stable"] = "stable";
            Stability["External"] = "external";
        })(Stability = exports.Stability || (exports.Stability = {}));
        var CollectionKind;
        (function(CollectionKind) {
            CollectionKind["Array"] = "array";
            CollectionKind["Map"] = "map";
        })(CollectionKind = exports.CollectionKind || (exports.CollectionKind = {}));
        var PrimitiveType;
        (function(PrimitiveType) {
            PrimitiveType["Date"] = "date";
            PrimitiveType["String"] = "string";
            PrimitiveType["Number"] = "number";
            PrimitiveType["Boolean"] = "boolean";
            PrimitiveType["Json"] = "json";
            PrimitiveType["Any"] = "any";
        })(PrimitiveType = exports.PrimitiveType || (exports.PrimitiveType = {}));
        exports.CANONICAL_ANY = {
            primitive: PrimitiveType.Any
        };
        function isNamedTypeReference(ref) {
            var _a;
            return !!((_a = ref) === null || _a === void 0 ? void 0 : _a.fqn);
        }
        exports.isNamedTypeReference = isNamedTypeReference;
        function isPrimitiveTypeReference(ref) {
            var _a;
            return !!((_a = ref) === null || _a === void 0 ? void 0 : _a.primitive);
        }
        exports.isPrimitiveTypeReference = isPrimitiveTypeReference;
        function isCollectionTypeReference(ref) {
            var _a;
            return !!((_a = ref) === null || _a === void 0 ? void 0 : _a.collection);
        }
        exports.isCollectionTypeReference = isCollectionTypeReference;
        function isUnionTypeReference(ref) {
            var _a;
            return !!((_a = ref) === null || _a === void 0 ? void 0 : _a.union);
        }
        exports.isUnionTypeReference = isUnionTypeReference;
        function isMethod(callable) {
            return !!callable.name;
        }
        exports.isMethod = isMethod;
        var TypeKind;
        (function(TypeKind) {
            TypeKind["Class"] = "class";
            TypeKind["Enum"] = "enum";
            TypeKind["Interface"] = "interface";
        })(TypeKind = exports.TypeKind || (exports.TypeKind = {}));
        function isClassType(type) {
            return (type === null || type === void 0 ? void 0 : type.kind) === TypeKind.Class;
        }
        exports.isClassType = isClassType;
        function isInterfaceType(type) {
            return (type === null || type === void 0 ? void 0 : type.kind) === TypeKind.Interface;
        }
        exports.isInterfaceType = isInterfaceType;
        function isEnumType(type) {
            return (type === null || type === void 0 ? void 0 : type.kind) === TypeKind.Enum;
        }
        exports.isEnumType = isEnumType;
        function isClassOrInterfaceType(type) {
            return isClassType(type) || isInterfaceType(type);
        }
        exports.isClassOrInterfaceType = isClassOrInterfaceType;
        function describeTypeReference(type) {
            if (type === undefined) {
                return "void";
            }
            if (isNamedTypeReference(type)) {
                return type.fqn;
            }
            if (isPrimitiveTypeReference(type)) {
                return type.primitive;
            }
            if (isCollectionTypeReference(type)) {
                return `${type.collection.kind}<${describeTypeReference(type.collection.elementtype)}>`;
            }
            if (isUnionTypeReference(type)) {
                const unionType = type.union.types.map(describeTypeReference).join(" | ");
                return unionType;
            }
            throw new Error("Unrecognized type reference");
        }
        exports.describeTypeReference = describeTypeReference;
        function isDeprecated(entity) {
            var _a;
            return ((_a = entity === null || entity === void 0 ? void 0 : entity.docs) === null || _a === void 0 ? void 0 : _a.stability) === Stability.Deprecated;
        }
        exports.isDeprecated = isDeprecated;
    },
    5585: (__unused_webpack_module, exports) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
    },
    1804: function(__unused_webpack_module, exports, __webpack_require__) {
        "use strict";
        var __createBinding = this && this.__createBinding || (Object.create ? function(o, m, k, k2) {
            if (k2 === undefined) k2 = k;
            Object.defineProperty(o, k2, {
                enumerable: true,
                get: function() {
                    return m[k];
                }
            });
        } : function(o, m, k, k2) {
            if (k2 === undefined) k2 = k;
            o[k2] = m[k];
        });
        var __exportStar = this && this.__exportStar || function(m, exports) {
            for (var p in m) if (p !== "default" && !exports.hasOwnProperty(p)) __createBinding(exports, m, p);
        };
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        __exportStar(__webpack_require__(2752), exports);
        __exportStar(__webpack_require__(5585), exports);
        __exportStar(__webpack_require__(1485), exports);
        __exportStar(__webpack_require__(5907), exports);
    },
    1485: (__unused_webpack_module, exports) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.NameTree = void 0;
        class NameTree {
            constructor() {
                this._children = {};
            }
            static of(assm) {
                var _a;
                const nameTree = new NameTree;
                for (const type of Object.values((_a = assm.types) !== null && _a !== void 0 ? _a : {})) {
                    nameTree.register(type.fqn);
                }
                return nameTree;
            }
            get children() {
                return this._children;
            }
            get fqn() {
                return this._fqn;
            }
            register(fqn, path = fqn.split(".")) {
                if (path.length === 0) {
                    this._fqn = fqn;
                } else {
                    const [head, ...rest] = path;
                    if (!this._children[head]) {
                        this._children[head] = new NameTree;
                    }
                    this._children[head].register(fqn, rest);
                }
                return this;
            }
        }
        exports.NameTree = NameTree;
    },
    5907: (__unused_webpack_module, exports, __webpack_require__) => {
        "use strict";
        Object.defineProperty(exports, "__esModule", {
            value: true
        });
        exports.validateAssembly = exports.schema = void 0;
        const jsonschema_1 = __webpack_require__(6872);
        exports.schema = __webpack_require__(373);
        function validateAssembly(obj) {
            const validator = new jsonschema_1.Validator;
            validator.addSchema(exports.schema);
            const result = validator.validate(obj, exports.schema, {
                nestedErrors: true
            });
            if (result.valid) {
                return obj;
            }
            throw new Error(`Invalid assembly:\n${result.toString()}`);
        }
        exports.validateAssembly = validateAssembly;
    },
    306: module => {
        "use strict";
        module.exports = JSON.parse('{"name":"@jsii/runtime","version":"1.30.0","description":"jsii runtime kernel process","license":"Apache-2.0","author":{"name":"Amazon Web Services","url":"https://aws.amazon.com"},"homepage":"https://github.com/aws/jsii","bugs":{"url":"https://github.com/aws/jsii/issues"},"repository":{"type":"git","url":"https://github.com/aws/jsii.git","directory":"packages/@jsii/runtime"},"engines":{"node":">= 10.3.0"},"main":"lib/index.js","types":"lib/index.d.ts","bin":{"jsii-runtime":"bin/jsii-runtime"},"scripts":{"build":"tsc --build && chmod +x bin/jsii-runtime && npx webpack-cli && npm run lint","watch":"tsc --build -w","lint":"eslint . --ext .js,.ts --ignore-path=.gitignore --ignore-pattern=webpack.config.js","lint:fix":"yarn lint --fix","test":"jest","test:update":"jest -u","package":"package-js"},"dependencies":{"@jsii/kernel":"^1.30.0","@jsii/spec":"^1.30.0"},"devDependencies":{"@scope/jsii-calc-base":"^1.30.0","@scope/jsii-calc-lib":"^1.30.0","@types/jest":"^26.0.23","@types/node":"^10.17.59","eslint":"^7.26.0","jest":"^26.6.3","jsii-build-tools":"^1.30.0","jsii-calc":"^3.20.120","prettier":"^2.3.0","source-map-loader":"^2.0.1","ts-jest":"^26.5.6","typescript":"~3.9.9","webpack":"^5.37.0","webpack-cli":"^4.7.0"}}');
    },
    373: module => {
        "use strict";
        module.exports = JSON.parse('{"$ref":"#/definitions/Assembly","$schema":"http://json-schema.org/draft-07/schema#","definitions":{"Assembly":{"description":"A JSII assembly specification.","properties":{"author":{"$ref":"#/definitions/Person","description":"The main author of this package."},"bin":{"additionalProperties":{"type":"string"},"default":"none","description":"List of bin-scripts","type":"object"},"bundled":{"additionalProperties":{"type":"string"},"default":"none","description":"List if bundled dependencies (these are not expected to be jsii\\nassemblies).","type":"object"},"contributors":{"default":"none","description":"Additional contributors to this package.","items":{"$ref":"#/definitions/Person"},"type":"array"},"dependencies":{"additionalProperties":{"type":"string"},"default":"none","description":"Direct dependencies on other assemblies (with semver), the key is the JSII\\nassembly name, and the value is a SemVer expression.","type":"object"},"dependencyClosure":{"additionalProperties":{"$ref":"#/definitions/AssemblyConfiguration"},"default":"none","description":"Target configuration for all the assemblies that are direct or transitive\\ndependencies of this assembly. This is needed to generate correct native\\ntype names for any transitively inherited member, in certain languages.","type":"object"},"description":{"description":"Description of the assembly, maps to \\"description\\" from package.json\\nThis is required since some package managers (like Maven) require it.","type":"string"},"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"fingerprint":{"description":"A fingerprint that can be used to determine if the specification has\\nchanged.","minLength":1,"type":"string"},"homepage":{"description":"The url to the project homepage. Maps to \\"homepage\\" from package.json.","type":"string"},"jsiiVersion":{"description":"The version of the jsii compiler that was used to produce this Assembly.","minLength":1,"type":"string"},"keywords":{"description":"Keywords that help discover or identify this packages with respects to it\'s\\nintended usage, audience, etc... Where possible, this will be rendered in\\nthe corresponding metadata section of idiomatic package manifests, for\\nexample NuGet package tags.","items":{"type":"string"},"type":"array"},"license":{"description":"The SPDX name of the license this assembly is distributed on.","type":"string"},"metadata":{"additionalProperties":{},"default":"none","description":"Arbitrary key-value pairs of metadata, which the maintainer chose to\\ndocument with the assembly. These entries do not carry normative\\nsemantics and their interpretation is up to the assembly maintainer.","type":"object"},"name":{"description":"The name of the assembly","minLength":1,"type":"string"},"readme":{"$ref":"#/definitions/ReadMe","default":"none","description":"The readme document for this module (if any)."},"repository":{"description":"The module repository, maps to \\"repository\\" from package.json\\nThis is required since some package managers (like Maven) require it.","properties":{"directory":{"default":"the root of the repository","description":"If the package is not in the root directory (for example, when part\\nof a monorepo), you should specify the directory in which it lives.","type":"string"},"type":{"description":"The type of the repository (``git``, ``svn``, ...)","type":"string"},"url":{"description":"The URL of the repository.","type":"string"}},"required":["type","url"],"type":"object"},"schema":{"description":"The version of the spec schema","enum":["jsii/0.10.0"],"type":"string"},"submodules":{"additionalProperties":{"allOf":[{"$ref":"#/definitions/SourceLocatable"},{"$ref":"#/definitions/Targetable"}],"description":"A submodule\\n\\nThe difference between a top-level module (the assembly) and a submodule is\\nthat the submodule is annotated with its location in the repository."},"default":"none","description":"Submodules declared in this assembly.","type":"object"},"targets":{"$ref":"#/definitions/AssemblyTargets","default":"none","description":"A map of target name to configuration, which is used when generating\\npackages for various languages."},"types":{"additionalProperties":{"anyOf":[{"allOf":[{"$ref":"#/definitions/TypeBase"},{"$ref":"#/definitions/ClassType"}]},{"allOf":[{"$ref":"#/definitions/TypeBase"},{"$ref":"#/definitions/EnumType"}]},{"allOf":[{"$ref":"#/definitions/TypeBase"},{"$ref":"#/definitions/InterfaceType"}]}],"description":"Represents a type definition (not a type reference)."},"default":"none","description":"All types in the assembly, keyed by their fully-qualified-name","type":"object"},"version":{"description":"The version of the assembly","minLength":1,"type":"string"}},"required":["author","description","fingerprint","homepage","jsiiVersion","license","name","repository","schema","version"],"type":"object"},"AssemblyConfiguration":{"description":"Shareable configuration of a jsii Assembly.","properties":{"readme":{"$ref":"#/definitions/ReadMe","default":"none","description":"The readme document for this module (if any)."},"submodules":{"additionalProperties":{"allOf":[{"$ref":"#/definitions/SourceLocatable"},{"$ref":"#/definitions/Targetable"}],"description":"A submodule\\n\\nThe difference between a top-level module (the assembly) and a submodule is\\nthat the submodule is annotated with its location in the repository."},"default":"none","description":"Submodules declared in this assembly.","type":"object"},"targets":{"$ref":"#/definitions/AssemblyTargets","default":"none","description":"A map of target name to configuration, which is used when generating\\npackages for various languages."}},"type":"object"},"AssemblyTargets":{"additionalProperties":{"additionalProperties":{},"type":"object"},"description":"Configurable targets for an asembly.","type":"object"},"Callable":{"description":"An Initializer or a Method.","properties":{"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"locationInModule":{"$ref":"#/definitions/SourceLocation","default":"none","description":"Where in the module this definition was found\\n\\nWhy is this not `locationInAssembly`? Because the assembly is the JSII\\nfile combining compiled code and its manifest, whereas this is referring\\nto the location of the source in the module the assembly was built from."},"overrides":{"default":"this member is not overriding anything","description":"The FQN of the parent type (class or interface) that this entity\\noverrides or implements. If undefined, then this entity is the first in\\nit\'s hierarchy to declare this entity.","type":"string"},"parameters":{"default":"none","description":"The parameters of the Initializer or Method.","items":{"$ref":"#/definitions/Parameter"},"type":"array"},"protected":{"default":false,"description":"Indicates if this Initializer or Method is protected (otherwise it is\\npublic, since private members are not modeled).","type":"boolean"},"variadic":{"default":false,"description":"Indicates whether this Initializer or Method is variadic or not. When\\n``true``, the last element of ``#parameters`` will also be flagged\\n``#variadic``.","type":"boolean"}},"type":"object"},"ClassType":{"description":"Represents classes.","properties":{"abstract":{"default":false,"description":"Indicates if this class is an abstract class.","type":"boolean"},"assembly":{"description":"The name of the assembly the type belongs to.","minLength":1,"type":"string"},"base":{"default":"no base class","description":"The FQN of the base class of this class, if it has one.","type":"string"},"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"fqn":{"description":"The fully qualified name of the type (``<assembly>.<namespace>.<name>``)","minLength":3,"type":"string"},"initializer":{"$ref":"#/definitions/Callable","default":"no initializer","description":"Initializer (constructor) method."},"interfaces":{"default":"none","description":"The FQNs of the interfaces this class implements, if any.","items":{"type":"string"},"type":"array","uniqueItems":true},"kind":{"description":"The kind of the type.","enum":["class"],"type":"string"},"locationInModule":{"$ref":"#/definitions/SourceLocation","default":"none","description":"Where in the module this definition was found\\n\\nWhy is this not `locationInAssembly`? Because the assembly is the JSII\\nfile combining compiled code and its manifest, whereas this is referring\\nto the location of the source in the module the assembly was built from."},"methods":{"default":"none","description":"List of methods.","items":{"$ref":"#/definitions/Method"},"type":"array"},"name":{"description":"The simple name of the type (MyClass).","minLength":1,"type":"string"},"namespace":{"default":"none","description":"The namespace of the type (`foo.bar.baz`).\\n\\nWhen undefined, the type is located at the root of the assembly (its\\n`fqn` would be like `<assembly>.<name>`).\\n\\nFor types inside other types or inside submodules, the `<namespace>` corresponds to\\nthe namespace-qualified name of the container (can contain multiple segments like:\\n`<ns1>.<ns2>.<ns3>`).\\n\\nIn all cases:\\n\\n <fqn> = <assembly>[.<namespace>].<name>","type":"string"},"properties":{"default":"none","description":"List of properties.","items":{"$ref":"#/definitions/Property"},"type":"array"}},"required":["assembly","fqn","kind","name"],"type":"object"},"CollectionKind":{"description":"Kinds of collections.","enum":["array","map"],"type":"string"},"CollectionTypeReference":{"description":"Reference to a collection type.","properties":{"collection":{"properties":{"elementtype":{"$ref":"#/definitions/TypeReference","description":"The type of an element (map keys are always strings)."},"kind":{"$ref":"#/definitions/CollectionKind","description":"The kind of collection."}},"required":["elementtype","kind"],"type":"object"}},"required":["collection"],"type":"object"},"Docs":{"description":"Key value pairs of documentation nodes.\\nBased on TSDoc.","properties":{"custom":{"additionalProperties":{"type":"string"},"default":"none","description":"Custom tags that are not any of the default ones","type":"object"},"default":{"default":"none","description":"Description of the default","type":"string"},"deprecated":{"default":"none","description":"If present, this block indicates that an API item is no longer supported\\nand may be removed in a future release.  The `@deprecated` tag must be\\nfollowed by a sentence describing the recommended alternative.\\nDeprecation recursively applies to members of a container. For example,\\nif a class is deprecated, then so are all of its members.","type":"string"},"example":{"default":"none","description":"Example showing the usage of this API item\\n\\nStarts off in running text mode, may switch to code using fenced code\\nblocks.","type":"string"},"remarks":{"default":"none","description":"Detailed information about an API item.\\n\\nEither the explicitly tagged `@remarks` section, otherwise everything\\npast the first paragraph if there is no `@remarks` tag.","type":"string"},"returns":{"default":"none","description":"The `@returns` block for this doc comment, or undefined if there is not\\none.","type":"string"},"see":{"default":"none","description":"A `@see` link with more information","type":"string"},"stability":{"description":"Whether the API item is beta/experimental quality","enum":["deprecated","experimental","external","stable"],"type":"string"},"subclassable":{"default":false,"description":"Whether this class or interface was intended to be subclassed/implemented\\nby library users.\\n\\nClasses intended for subclassing, and interfaces intended to be\\nimplemented by consumers, are held to stricter standards of API\\ncompatibility.","type":"boolean"},"summary":{"default":"none","description":"Summary documentation for an API item.\\n\\nThe first part of the documentation before hitting a `@remarks` tags, or\\nthe first line of the doc comment block if there is no `@remarks` tag.","type":"string"}},"type":"object"},"EnumMember":{"description":"Represents a member of an enum.","properties":{"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"name":{"description":"The name/symbol of the member.","type":"string"}},"required":["name"],"type":"object"},"EnumType":{"description":"Represents an enum type.","properties":{"assembly":{"description":"The name of the assembly the type belongs to.","minLength":1,"type":"string"},"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"fqn":{"description":"The fully qualified name of the type (``<assembly>.<namespace>.<name>``)","minLength":3,"type":"string"},"kind":{"description":"The kind of the type.","enum":["enum"],"type":"string"},"locationInModule":{"$ref":"#/definitions/SourceLocation","default":"none","description":"Where in the module this definition was found\\n\\nWhy is this not `locationInAssembly`? Because the assembly is the JSII\\nfile combining compiled code and its manifest, whereas this is referring\\nto the location of the source in the module the assembly was built from."},"members":{"description":"Members of the enum.","items":{"$ref":"#/definitions/EnumMember"},"type":"array"},"name":{"description":"The simple name of the type (MyClass).","minLength":1,"type":"string"},"namespace":{"default":"none","description":"The namespace of the type (`foo.bar.baz`).\\n\\nWhen undefined, the type is located at the root of the assembly (its\\n`fqn` would be like `<assembly>.<name>`).\\n\\nFor types inside other types or inside submodules, the `<namespace>` corresponds to\\nthe namespace-qualified name of the container (can contain multiple segments like:\\n`<ns1>.<ns2>.<ns3>`).\\n\\nIn all cases:\\n\\n <fqn> = <assembly>[.<namespace>].<name>","type":"string"}},"required":["assembly","fqn","kind","members","name"],"type":"object"},"InterfaceType":{"properties":{"assembly":{"description":"The name of the assembly the type belongs to.","minLength":1,"type":"string"},"datatype":{"default":false,"description":"True if this interface only contains properties. Different backends might\\nhave idiomatic ways to allow defining concrete instances such interfaces.\\nFor example, in Java, the generator will produce a PoJo and a builder\\nwhich will allow users to create a concrete object with data which\\nadheres to this interface.","type":"boolean"},"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"fqn":{"description":"The fully qualified name of the type (``<assembly>.<namespace>.<name>``)","minLength":3,"type":"string"},"interfaces":{"default":"none","description":"The FQNs of the interfaces this interface extends, if any.","items":{"type":"string"},"type":"array","uniqueItems":true},"kind":{"description":"The kind of the type.","enum":["interface"],"type":"string"},"locationInModule":{"$ref":"#/definitions/SourceLocation","default":"none","description":"Where in the module this definition was found\\n\\nWhy is this not `locationInAssembly`? Because the assembly is the JSII\\nfile combining compiled code and its manifest, whereas this is referring\\nto the location of the source in the module the assembly was built from."},"methods":{"default":"none","description":"List of methods.","items":{"$ref":"#/definitions/Method"},"type":"array"},"name":{"description":"The simple name of the type (MyClass).","minLength":1,"type":"string"},"namespace":{"default":"none","description":"The namespace of the type (`foo.bar.baz`).\\n\\nWhen undefined, the type is located at the root of the assembly (its\\n`fqn` would be like `<assembly>.<name>`).\\n\\nFor types inside other types or inside submodules, the `<namespace>` corresponds to\\nthe namespace-qualified name of the container (can contain multiple segments like:\\n`<ns1>.<ns2>.<ns3>`).\\n\\nIn all cases:\\n\\n <fqn> = <assembly>[.<namespace>].<name>","type":"string"},"properties":{"default":"none","description":"List of properties.","items":{"$ref":"#/definitions/Property"},"type":"array"}},"required":["assembly","fqn","kind","name"],"type":"object"},"Method":{"description":"A method with a name (i.e: not an initializer).","properties":{"abstract":{"default":false,"description":"Is this method an abstract method (this means the class will also be an abstract class)","type":"boolean"},"async":{"default":false,"description":"Indicates if this is an asynchronous method (it will return a promise).","type":"boolean"},"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"locationInModule":{"$ref":"#/definitions/SourceLocation","default":"none","description":"Where in the module this definition was found\\n\\nWhy is this not `locationInAssembly`? Because the assembly is the JSII\\nfile combining compiled code and its manifest, whereas this is referring\\nto the location of the source in the module the assembly was built from."},"name":{"description":"The name of the method. Undefined if this method is a initializer.","type":"string"},"overrides":{"default":"this member is not overriding anything","description":"The FQN of the parent type (class or interface) that this entity\\noverrides or implements. If undefined, then this entity is the first in\\nit\'s hierarchy to declare this entity.","type":"string"},"parameters":{"default":"none","description":"The parameters of the Initializer or Method.","items":{"$ref":"#/definitions/Parameter"},"type":"array"},"protected":{"default":false,"description":"Indicates if this Initializer or Method is protected (otherwise it is\\npublic, since private members are not modeled).","type":"boolean"},"returns":{"$ref":"#/definitions/OptionalValue","default":"void","description":"The return type of the method (`undefined` if `void`)"},"static":{"default":false,"description":"Indicates if this is a static method.","type":"boolean"},"variadic":{"default":false,"description":"Indicates whether this Initializer or Method is variadic or not. When\\n``true``, the last element of ``#parameters`` will also be flagged\\n``#variadic``.","type":"boolean"}},"required":["name"],"type":"object"},"NamedTypeReference":{"description":"Reference to a named type, defined by this assembly or one of its\\ndependencies.","properties":{"fqn":{"description":"The fully-qualified-name of the type (can be located in the\\n``spec.types[fqn]``` of the assembly that defines the type).","type":"string"}},"required":["fqn"],"type":"object"},"OptionalValue":{"description":"A value that can possibly be optional.","properties":{"optional":{"default":false,"description":"Determines whether the value is, indeed, optional.","type":"boolean"},"type":{"$ref":"#/definitions/TypeReference","description":"The declared type of the value, when it\'s present."}},"required":["type"],"type":"object"},"Parameter":{"description":"Represents a method parameter.","properties":{"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"name":{"description":"The name of the parameter.","minLength":1,"type":"string"},"optional":{"default":false,"description":"Determines whether the value is, indeed, optional.","type":"boolean"},"type":{"$ref":"#/definitions/TypeReference","description":"The declared type of the value, when it\'s present."},"variadic":{"default":false,"description":"Whether this is the last parameter of a variadic method. In such cases,\\nthe `#type` attribute is the type of each individual item of the variadic\\narguments list (as opposed to some array type, as for example TypeScript\\nwould model it).","type":"boolean"}},"required":["name","type"],"type":"object"},"Person":{"description":"Metadata about people or organizations associated with the project that\\nresulted in the Assembly. Some of this metadata is required in order to\\npublish to certain package repositories (for example, Maven Central), but is\\nnot normalized, and the meaning of fields (role, for example), is up to each\\nproject maintainer.","properties":{"email":{"default":"none","description":"The email of the person","type":"string"},"name":{"description":"The name of the person","type":"string"},"organization":{"default":false,"description":"If true, this person is, in fact, an organization","type":"boolean"},"roles":{"description":"A list of roles this person has in the project, for example `maintainer`,\\n`contributor`, `owner`, ...","items":{"type":"string"},"type":"array"},"url":{"default":"none","description":"The URL for the person","type":"string"}},"required":["name","roles"],"type":"object"},"PrimitiveType":{"description":"Kinds of primitive types.","enum":["any","boolean","date","json","number","string"],"type":"string"},"PrimitiveTypeReference":{"description":"Reference to a primitive type.","properties":{"primitive":{"$ref":"#/definitions/PrimitiveType","description":"If this is a reference to a primitive type, this will include the\\nprimitive type kind."}},"required":["primitive"],"type":"object"},"Property":{"description":"A class property.","properties":{"abstract":{"default":false,"description":"Indicates if this property is abstract","type":"boolean"},"const":{"default":false,"description":"A hint that indicates that this static, immutable property is initialized\\nduring startup. This allows emitting \\"const\\" idioms in different target\\nlanguages. Implies `static` and `immutable`.","type":"boolean"},"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"immutable":{"default":false,"description":"Indicates if this property only has a getter (immutable).","type":"boolean"},"locationInModule":{"$ref":"#/definitions/SourceLocation","default":"none","description":"Where in the module this definition was found\\n\\nWhy is this not `locationInAssembly`? Because the assembly is the JSII\\nfile combining compiled code and its manifest, whereas this is referring\\nto the location of the source in the module the assembly was built from."},"name":{"description":"The name of the property.","minLength":1,"type":"string"},"optional":{"default":false,"description":"Determines whether the value is, indeed, optional.","type":"boolean"},"overrides":{"default":"this member is not overriding anything","description":"The FQN of the parent type (class or interface) that this entity\\noverrides or implements. If undefined, then this entity is the first in\\nit\'s hierarchy to declare this entity.","type":"string"},"protected":{"default":false,"description":"Indicates if this property is protected (otherwise it is public)","type":"boolean"},"static":{"default":false,"description":"Indicates if this is a static property.","type":"boolean"},"type":{"$ref":"#/definitions/TypeReference","description":"The declared type of the value, when it\'s present."}},"required":["name","type"],"type":"object"},"ReadMe":{"description":"README information","properties":{"markdown":{"type":"string"}},"required":["markdown"],"type":"object"},"SourceLocatable":{"description":"Indicates that an entity has a source location","properties":{"locationInModule":{"$ref":"#/definitions/SourceLocation","default":"none","description":"Where in the module this definition was found\\n\\nWhy is this not `locationInAssembly`? Because the assembly is the JSII\\nfile combining compiled code and its manifest, whereas this is referring\\nto the location of the source in the module the assembly was built from."}},"type":"object"},"SourceLocation":{"description":"Where in the module source the definition for this API item was found","properties":{"filename":{"description":"Relative filename","type":"string"},"line":{"description":"1-based line number in the indicated file","type":"number"}},"required":["filename","line"],"type":"object"},"Targetable":{"description":"A targetable module-like thing\\n\\nHas targets and a readme. Used for Assemblies and Submodules.","properties":{"readme":{"$ref":"#/definitions/ReadMe","default":"none","description":"The readme document for this module (if any)."},"targets":{"$ref":"#/definitions/AssemblyTargets","default":"none","description":"A map of target name to configuration, which is used when generating\\npackages for various languages."}},"type":"object"},"TypeBase":{"description":"Common attributes of a type definition.","properties":{"assembly":{"description":"The name of the assembly the type belongs to.","minLength":1,"type":"string"},"docs":{"$ref":"#/definitions/Docs","default":"none","description":"Documentation for this entity."},"fqn":{"description":"The fully qualified name of the type (``<assembly>.<namespace>.<name>``)","minLength":3,"type":"string"},"kind":{"$ref":"#/definitions/TypeKind","description":"The kind of the type."},"locationInModule":{"$ref":"#/definitions/SourceLocation","default":"none","description":"Where in the module this definition was found\\n\\nWhy is this not `locationInAssembly`? Because the assembly is the JSII\\nfile combining compiled code and its manifest, whereas this is referring\\nto the location of the source in the module the assembly was built from."},"name":{"description":"The simple name of the type (MyClass).","minLength":1,"type":"string"},"namespace":{"default":"none","description":"The namespace of the type (`foo.bar.baz`).\\n\\nWhen undefined, the type is located at the root of the assembly (its\\n`fqn` would be like `<assembly>.<name>`).\\n\\nFor types inside other types or inside submodules, the `<namespace>` corresponds to\\nthe namespace-qualified name of the container (can contain multiple segments like:\\n`<ns1>.<ns2>.<ns3>`).\\n\\nIn all cases:\\n\\n <fqn> = <assembly>[.<namespace>].<name>","type":"string"}},"required":["assembly","fqn","kind","name"],"type":"object"},"TypeKind":{"description":"Kinds of types.","enum":["class","enum","interface"],"type":"string"},"TypeReference":{"anyOf":[{"$ref":"#/definitions/NamedTypeReference"},{"$ref":"#/definitions/PrimitiveTypeReference"},{"$ref":"#/definitions/CollectionTypeReference"},{"$ref":"#/definitions/UnionTypeReference"}],"description":"A reference to a type (primitive, collection or fqn)."},"UnionTypeReference":{"description":"Reference to a union type.","properties":{"union":{"description":"Indicates that this is a union type, which means it can be one of a set\\nof types.","properties":{"types":{"description":"All the possible types (including the primary type).","items":{"$ref":"#/definitions/TypeReference"},"minItems":2,"type":"array"}},"required":["types"],"type":"object"}},"required":["union"],"type":"object"}}}');
    },
    2357: module => {
        "use strict";
        module.exports = require("assert");
    },
    4293: module => {
        "use strict";
        module.exports = require("buffer");
    },
    3129: module => {
        "use strict";
        module.exports = require("child_process");
    },
    7619: module => {
        "use strict";
        module.exports = require("constants");
    },
    6417: module => {
        "use strict";
        module.exports = require("crypto");
    },
    8614: module => {
        "use strict";
        module.exports = require("events");
    },
    5747: module => {
        "use strict";
        module.exports = require("fs");
    },
    2282: module => {
        "use strict";
        module.exports = require("module");
    },
    2087: module => {
        "use strict";
        module.exports = require("os");
    },
    5622: module => {
        "use strict";
        module.exports = require("path");
    },
    1765: module => {
        "use strict";
        module.exports = require("process");
    },
    2413: module => {
        "use strict";
        module.exports = require("stream");
    },
    4304: module => {
        "use strict";
        module.exports = require("string_decoder");
    },
    8835: module => {
        "use strict";
        module.exports = require("url");
    },
    1669: module => {
        "use strict";
        module.exports = require("util");
    },
    2184: module => {
        "use strict";
        module.exports = require("vm");
    },
    8761: module => {
        "use strict";
        module.exports = require("zlib");
    }
};

var __webpack_module_cache__ = {};

function __webpack_require__(moduleId) {
    var cachedModule = __webpack_module_cache__[moduleId];
    if (cachedModule !== undefined) {
        return cachedModule.exports;
    }
    var module = __webpack_module_cache__[moduleId] = {
        id: moduleId,
        loaded: false,
        exports: {}
    };
    __webpack_modules__[moduleId].call(module.exports, module, module.exports, __webpack_require__);
    module.loaded = true;
    return module.exports;
}

(() => {
    __webpack_require__.nmd = module => {
        module.paths = [];
        if (!module.children) module.children = [];
        return module;
    };
})();

var __webpack_exports__ = {};

(() => {
    "use strict";
    var exports = __webpack_exports__;
    var __webpack_unused_export__;
    var _a;
    __webpack_unused_export__ = {
        value: true
    };
    const process = __webpack_require__(1765);
    const packageInfo = __webpack_require__(306);
    const host_1 = __webpack_require__(7905);
    const in_out_1 = __webpack_require__(6156);
    const sync_stdio_1 = __webpack_require__(1416);
    const name = packageInfo.name;
    const version = packageInfo.version;
    const noStack = !!process.env.JSII_NOSTACK;
    const debug = !!process.env.JSII_DEBUG;
    const stdio = new sync_stdio_1.SyncStdio({
        errorFD: (_a = process.stderr.fd) !== null && _a !== void 0 ? _a : 2,
        readFD: 3,
        writeFD: 3
    });
    const inout = new in_out_1.InputOutput(stdio);
    const host = new host_1.KernelHost(inout, {
        debug,
        noStack
    });
    host.once("exit", process.exit.bind(process));
    inout.write({
        hello: `${name}@${version}`
    });
    inout.debug = debug;
    host.run();
})();
//# sourceMappingURL=program.js.map