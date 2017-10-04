# MongoDB for Quilt

This repository implements a MongoDB specification for Quilt.js.  See
[Quilt](http://quilt.io) for more information.

```javascript
var Mongo = require("github.com/quilt/mongo");

var deployment = createDeployment({});
deployment.deploy(new Mongo(3));
```
