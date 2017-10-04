const { Container, allow } = require('@quilt/quilt');

const image = 'quilt/mongo';

function getHostname(c) {
  return c.getHostname();
}

function Mongo(nWorker) {
    this.cluster = []
    for (i = 0; i < nWorker; i++)
    {
	this.cluster.push(new Container('mongo', image));
    }

  const hostnames = this.cluster.map(getHostname).join(',');
  this.cluster.forEach((m) => {
    m.setEnv('MEMBERS', hostnames);
  });

  // The initiator is choosen completley arbitrarily.
  this.cluster[0].setEnv('INITIATOR', 'true');

  allow(this.cluster, this.cluster, this.port);
}

Mongo.prototype.port = 27017;

Mongo.prototype.deploy = function deploy(deployment) {
  deployment.deploy(this.cluster);
};

Mongo.prototype.allowFrom = function allowFrom(from, p) {
  allow(from, this.cluster, p);
};

Mongo.prototype.uri = function uri(dbName) {
  return `mongodb://${this.cluster.map(getHostname).join(',')}/${dbName}`;
};

module.exports = Mongo;
