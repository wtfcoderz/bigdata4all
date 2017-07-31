
//Require the dev-dependencies
let chai = require('chai');
let chaiHttp = require('chai-http');
let should = chai.should();

chai.use(chaiHttp);

describe('/GET health', () => {
 it('it should return 200', (done) => {
   chai.request(process.env.API_URL)
   .get('/health')
   .end((err, res) => {
     res.should.have.status(200);
     done();
    });
  });
});
