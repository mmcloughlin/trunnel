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
	find ${REF}/examples ${REF}/test/valid -name '*.trunnel' > $@

.PHONY: FORCE
