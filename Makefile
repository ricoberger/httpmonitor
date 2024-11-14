.PHONY: release-major
release-major:
	$(eval MAJORVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print "v"$$1+1".0.0"}'))
	@git checkout main
	@git pull
	@git tag -a $(MAJORVERSION) -m 'Release $(MAJORVERSION)'
	@git push origin --tags

.PHONY: release-minor
release-minor:
	$(eval MINORVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print "v"$$1"."$$2+1".0"}'))
	@git checkout main
	@git pull
	@git tag -a $(MINORVERSION) -m 'Release $(MINORVERSION)'
	@git push origin --tags

.PHONY: release-patch
release-patch:
	$(eval PATCHVERSION=$(shell git describe --tags --abbrev=0 | sed s/v// | awk -F. '{print "v"$$1"."$$2"."$$3+1}'))
	@git checkout main
	@git pull
	@git tag -a $(PATCHVERSION) -m 'Release $(PATCHVERSION)'
	@git push origin --tags
