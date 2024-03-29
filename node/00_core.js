// Copyright 2021 @abdfnx. All rights reserved. Apache 2.0 license.
// Define op codes
const __ops = {
  FSWrite:     1,
  FSRead:      2,
  FSExists:    3,
  FSDirExists: 4,
  FSCwd:       5,
  FSStat:      6,
  FSRemove:    7,
  FSMkdir:     9,
  Log:        10,
  Env:        11,
  Walk:       14,
  Plugin:     15,
  Fetch:      20,
  Serve:      25,
};

((window) => {
  let initialized = false;
  let promiseTable = {};
  let promiseNextId = 1;

  function init() {
    if (initialized) return;
    initialized = true;
    globalThis.__recv(__recvAsync);
  }

  function __recvAsync(id, val) {
    if (!id) return;
    return promiseTable[id].resolve(val);
  }

  async function __sendAsync(op, cb, ...args) {
    init();
    const id = promiseNextId++;

    if (typeof cb == "function") {
      promiseTable[id] = { resolve: cb };
      globalThis.__send(op, ...[id, ...args]);
    } else {
      let resolve, reject;
      const promise = new Promise((resolve_, reject_) => {
        resolve = resolve_;
        reject = reject_;
      });

      promise.resolve = resolve;
      promise.reject = reject;

      promiseTable[id] = promise;

      globalThis.__send(op, ...[id, ...args]);

      const res = await promise;

      if (res.ok) return res.ok;
      else if (res.err) return res.err;
      else throw new Error("Unknown error");
    }
  }

  Object.assign(window, {
    __sendAsync,
  });
})(globalThis);
