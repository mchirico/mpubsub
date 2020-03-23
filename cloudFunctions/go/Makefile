

log:
	~/google-cloud-sdk/bin/gcloud functions logs read --limit 50

deploy:
	~/google-cloud-sdk/bin/gcloud functions deploy GoPubSub --runtime go113 --trigger-topic gocloud

delete:
	~/google-cloud-sdk/bin/gcloud functions delete GoPubSub
