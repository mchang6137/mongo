const { createDeployment, Machine, Range, githubKeys } = require('@quilt/quilt');
const Mongo = require('./mongo.js');

const nWorker = 3;

const deployment = createDeployment({});

const baseMachine = new Machine({
  provider: 'Amazon',
  cpu: new Range(2),
  ram: new Range(2),
  sshKeys: githubKeys('ejj'),
});
const mongo = new Mongo(nWorker);

deployment.deploy(baseMachine.asMaster());
deployment.deploy(baseMachine.asWorker().replicate(nWorker));
deployment.deploy(mongo);
