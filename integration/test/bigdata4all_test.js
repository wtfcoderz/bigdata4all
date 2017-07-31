let chai = require('chai');
let chaiHttp = require('chai-http');
let should = chai.should();
let assert = require('assert');

var api_url = 'http://api';

chai.use(chaiHttp);

describe('/GET health', () => {
  it('it should return 200', (done) => {
    chai.request(api_url)
    .get('/health')
    .end((err, res) => {
      if (err) assert(false, 'Error in /health');
      res.should.have.status(200);
      done();
    });
  });
});
