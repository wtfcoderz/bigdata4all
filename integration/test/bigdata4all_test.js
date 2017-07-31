
//Require the dev-dependencies
let chai = require('chai');
let chaiHttp = require('chai-http');
let should = chai.should();

var api_url = process.env.API_URL || 'http://api';

chai.use(chaiHttp);

describe('/GET health', () => {
 it('it should return 200', (done) => {
   chai.request(api_url)
   .get('/health')
   .end((err, res) => {
     res.should.have.status(200);
     done();
    });
  });
});
