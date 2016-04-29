 


 # CORE FUNCTIONS: ROUTETABLES


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.routetables.scan](#scan)
 > * [aws.routetables.describeForSubnet](#describeForSubnet)
 > * [aws.routetables.describe](#describe)
 > * [aws.routetables.create](#create)
 > * [aws.routetables.delete](#delete)
 > * [aws.routetables.associate](#associate)
 > * [aws.routetables.disassociate](#disassociate)
 > * [aws.routetables.deleteAssociation](#deleteAssociation)

 This API allows resource handlers to manipulate routing tables for subnets.

 ## AWS.ROUTETABLES.SCAN
 <a name="scan"></a>
 `aws.routetables.scan(region);`

 Query routetables.

 Example:

 ```

 var tables = aws.routetables.scan("us-east-1");

 ```

 ## AWS.ROUTETABLES.DESCRIBEFORSUBNET
 <a name="describeForSubnet"></a>
 `aws.routetables.describeForSubnet(region, subnet-id);`

 Get routetables associated with the supplied subnet.

 Example:

 ```

 var tables = aws.routetables.describeForSubnet("us-east-1", "subnet-abc");

 ```

 ## AWS.ROUTETABLES.DESCRIBE
 <a name="describe"></a>
 `aws.routetables.describe(region, route-table-id);`

 Get info about the supplied route table.

 Example:

 ```

 var table = aws.routetables.describe("us-east-1", "routetable-abc");

 ```

 ## AWS.ROUTETABLES.CREATE
 <a name="create"></a>
 `aws.routetables.create(region, vpc-id);`

 Create a route table in a vpc.

 Example:

 ```

 var table = aws.routetables.create("us-east-1", "vpc-123");

 ```

 ## AWS.ROUTETABLES.DELETE
 <a name="delete"></a>
 `aws.routetables.delete(region, subnet-id);`

 Delete a route table.

 Example:

 ```

 aws.routetables.delete("us-east-1", "routetable-123");

 ```

 ## AWS.ROUTETABLES.ASSOCIATE
 <a name="associate"></a>
 `aws.routetables.associate(region, subnet-id, route-table-id);`

 Assocate a route table with a subnet.

 Example:

 ```

 aws.routetables.associate("us-east-1", "subnet-abc", "routetable-123");

 ```

 ## AWS.ROUTETABLES.DISASSOCIATE
 <a name="disassociate"></a>
 `aws.routetables.disassociate(region, association-id);`

 Disassocate a route table and a subnet.

 Example:

 ```

 aws.routetables.disassociate("us-east-1", "association-xyz");

 ```

 ## AWS.ROUTETABLES.DELETEASSOCIATION
 <a name="deleteAssociation"></a>
 `aws.routetables.deleteAssociation(region, association-id);`

 Delete a route table association.

 Example:

 ```

 aws.routetables.deleteAssociation("us-east-1", "association-xyz");

 ```


