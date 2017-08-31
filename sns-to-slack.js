
// Based on https://gist.github.com/benyanke/862e446e5a816551928d8acc2d98b752
// Handles CloudWatch alerts via SNS as Slack Attachments, instead of plaintext.
//AWS Lambda function for forwarding SNS notifications to Slack

console.log('Loading function');

const https = require('https');
const url = require('url');
// to get the slack hook url, go into slack admin and create a new "Incoming Webhook" integration
const slack_url = 'https://hooks.slack.com/services/...'; // put your webhook URL here// Added by Ben Yanke
const slack_req_opts = url.parse(slack_url);
slack_req_opts.method = 'POST';
slack_req_opts.headers = {'Content-Type': 'application/json'};

exports.handler = function(event, context) {
  (event.Records || []).forEach(function (rec) {
    if (rec.Sns) {
      var req = https.request(slack_req_opts, function (res) {
        if (res.statusCode === 200) {
          context.succeed('posted to slack');
        } else {
          context.fail('status code: ' + res.statusCode);
        }
      });
      
      req.on('error', function(e) {
        console.log('problem with request: ' + e.message);
        context.fail(e.message);
      });
      
      
      
      // If event is a CloudWatch event
      if (rec.Sns.Subject.startsWith("ALARM:")) {

        cloudWatchMessage = JSON.parse(rec.Sns.Message)
        
          linkToInstance = "https://" + cloudWatchMessage.Region.toLowerCase() + 
                ".console.aws.amazon.com/ec2/v2/home?region=" + 
                cloudWatchMessage.Region.toLowerCase() + "#Instances:" + 
                cloudWatchMessage.Trigger.Dimensions[0].name + 
                "=" + cloudWatchMessage.Trigger.Dimensions[0].value + 
                ";sort=tag:Name";
                
         linkToAllInstances = "https://" + cloudWatchMessage.Region.toLowerCase() + 
                ".console.aws.amazon.com/ec2/v2/home?region=" + 
                cloudWatchMessage.Region.toLowerCase() + "#Instances:sort=tag:Name";
          
         linkToVpnConnections = "https://" + cloudWatchMessage.Region.toLowerCase() + 
                ".console.aws.amazon.com/vpc/home?region=" + 
                cloudWatchMessage.Region.toLowerCase() + "#VpnConnections";
          
         linkToConsole = "https://" + cloudWatchMessage.Region.toLowerCase() + 
                ".console.aws.amazon.com/console/home?region=" + 
                cloudWatchMessage.Region.toLowerCase();
          
          // Handle different types of errors 
          
          // EC2
          if(cloudWatchMessage.Trigger.Namespace == "AWS/EC2") {
            consoleLink = "<" + linkToInstance + "|Click here to open affected EC2 instance>";
            
          // VPN
          } else if(cloudWatchMessage.Trigger.Namespace == "AWS/VPN") {
            consoleLink = "<" + linkToVpnConnections + "|Click here to open VPN connections>";
            
          // All others
          } else {              
            consoleLink = "<" + linkToConsole + "|Click here to open AWS console>";  
          }
          
          
          req.write(JSON.stringify(
             {
                "attachments": [
                    {
                        "fallback": "CloudWatch Alert: " + cloudWatchMessage.AlarmName + " was triggered in " + cloudWatchMessage.Region,
                        "pretext": "New CloudWatch Alert",
                        "title": "CloudWatch Alarm: " + cloudWatchMessage.AlarmName,
                        "title_link" : linkToAllInstances,
                        "text": cloudWatchMessage.NewStateReason,
                        "color": "danger",
                        "author_name" : "Service: " + cloudWatchMessage.Trigger.Namespace,
                        "author_link" : linkToAllInstances,
                        "fields": [
                            {
                                "title": "Alarm",
                                "value": cloudWatchMessage.AlarmName,
                                "short": true
                            },
                            {
                                "title": "Alarm Description",
                                "value": cloudWatchMessage.AlarmDescription,
                                "short": true
                            },
                            {
                                "title": "Region",
                                "value": cloudWatchMessage.Region,
                                "short": true
                            },  
                            {
                                "title": "Environment",
                                "value": cloudWatchMessage.Trigger.Namespace,
                                "short": true
                            },  
                            {
                                "title": "AWS Console",
                                "value": consoleLink,
                                "short": true
                            }

                        ]
                    }
                ]
            }
          ))
          
      }
      
      // Otherwise, if not a CloudWatch event, send full alert subject and message
      else {

          req.write(JSON.stringify(
            {
                text: "*Subject:*   " + rec.Sns.Subject + "\n *Message:* " + rec.Sns.Message
            }
          ))
          
      }
      
      req.end();
    }
  });
};
