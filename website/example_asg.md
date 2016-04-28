  
 
 # Autoscaling Groups Example
 
 This is a more advanced example, demonstrating how to work with
 Mithras, which is an agentless system, in the context of AWS's
 Autoscaling Groups, which present an interesting context and
 challenge.
 
 If you're not using ASGs, then don't worry about this example.
 
 Usage:
 
     mithras -v run -f example/asg.js
 
 This example works with Mithas daemon mode.  After you've set up
 the ASG, using the above script, you'll run:
 
 
     mithras -v daemon start -f example/asg_daemon.js
 
 

