# lattigo and contract
[Chinese Simplified](./doc/README_ZH.md)
## Background
Enterprise blockchain is an area of great concern, and many enterprises or individuals are constantly promoting the application and practice of enterprise blockchain in China's government, justice, finance, and supply chain. In China, Hyperledger Fabric is a leader in enterprise-level blockchains, with a large number of business landings. At the same time, due to the government's emphasis on information security, enterprises have made more efforts in related matters, such as China's national commercial password, privacy computing and so on.  
What I want to share with you here today is related to private computing. China has an old saying called "killing chickens with cattle knives", figuratively doing small things without using large force, in the end it is cost control, which is also a point that enterprises are very concerned about, how to use relatively reasonable costs to complete the demand? Privacy computing is a big topic, but many times the business needs only a small part of it.  
lattigo is a lattice-based multi-party homomorphic encryption library implemented with go, which enables us to perform calculations on ciphertext, which does not expose private data, which is enough to meet most business scenarios. Because it is small enough to be flexibly applied in smart contracts, the combination achieves computable but immutable advantages.
## Practice
### Scenario
Company A needs to count the income and expenditure of this month to make financial statements, in order to avoid leaking the company's operating conditions, it is hoped that the income and expenditure details can be stored in ciphertext and relevant statistics can be carried out. For example:
|  Matter | Income and expenditure (yuan) |
|  ----  | ----  |
| Team building | -10000 |
| Party A settles the project amount | +1000000 |
| settle supplier payments |-50,000 |
| equipment procurement |-50000 |
| balance | 890,000 |  
### Contract
- CreateReport: Create a report
- SubmitData: Commit data
- QueryData: Query and count the results of the currently submitted data

For the above scenario, because the person in charge generates a set of key pairs and creates a report, the public key is handed over to the individual or group that needs to submit the data, and the submitter submits the data encrypted by the public key, the person in charge can query and count the results of the current submitted data after the data is submitted, decrypt the plaintext of the result with the private key, and of course, directly obtain the encrypted report and decrypt it.

## Insufficient
- It is not well connected to the existing fabric chaincode call system, and it needs to be packaged for services