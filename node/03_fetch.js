// Copyright 2021 @abdfnx. All rights reserved. Apache 2.0 license.

globalThis.fetch = async function (url) {
  return globalThis.__sendAsync(__ops.Fetch, false, url);
};
