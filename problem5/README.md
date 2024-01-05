## _Getting started:_

- Refer to `README.md` in `main` branch


## Consensus Breaking Change

### Description of the Change

*   **Change**: Addition of an `age` field (type: integer) to the `User` data structure.
*   **Reason for Change**: To include the age of users as part of their profile information in the blockchain.

### Impact of the Change

*   **Consensus-Breaking**: This change is consensus-breaking, meaning it will alter the way transactions are validated across the network.
*   **Validation Rules Change**: Existing nodes will reject transactions that include the new `age` field, as they are not aware of this new structure.
*   **Potential for Blockchain Fork**: If some nodes in the network adopt this change while others do not, it could lead to a fork in the blockchain. One version will have blocks with the `age` field, while the other will not.