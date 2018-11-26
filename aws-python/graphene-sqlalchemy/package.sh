rm -rf _package*
mkdir _package
cd _package && pip install -r ../requirements.txt --target .
zip -r9 ../_package.zip .
cd .. && zip -g _package.zip main.py
rm -rf _package