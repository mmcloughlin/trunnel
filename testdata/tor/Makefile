REF=../../ref/tor

all: README.md import

README.md: manifest.txt
	echo 'Test files found in tor repository version:' > $@
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
	find ${REF} -name '*.trunnel' > $@

.PHONY: FORCE
