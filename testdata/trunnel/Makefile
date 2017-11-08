REF=../../ref/trunnel

all: README.md import

README.md: manifest.txt
	echo 'Test files found in trunnel repository version:' > $@
	echo '```' >> $@
	git submodule status ${REF} >> $@
	echo '```' >> $@
	echo 'Original locations:' >> $@
	echo '```' >> $@
	cat $< >> $@
	echo '```' >> $@

import: manifest.txt
	xargs -i cp -v {} . < $<

manifest.txt: FORCE
	# excluding tor.trunnel because it contains examples of unions without
	# trailing semi-colons
	find ${REF}/examples ${REF}/test/valid -name '*.trunnel' | grep -v 'tor\.trunnel' > $@

.PHONY: FORCE
