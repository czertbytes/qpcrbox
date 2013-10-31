package main

import (
	"bytes"
	"archive/zip"
	"fmt"
)

type ODSExport struct {}

func (export *ODSExport) Export(e *Experiment) ([]byte, error) {
	content := new(bytes.Buffer)
	zw := zip.NewWriter(content)

	files := []FileData{
		FileData{"META-INF/manifest.xml", odsMetaInfManifestXmlFileContent()},
		FileData{"content.xml", odsContentXmlFileContent(e)},
		FileData{"meta.xml", odsMetaXmlFileContent()},
		FileData{"mimetype", odsMimetypeFileContent()},
		FileData{"settings.xml", odsSettingsXmlFileContent()},
		FileData{"styles.xml", odsStylesXmlFileContent()},
	}

	for _, fd := range files {
		addToArchive(zw, fd)
	}

	zw.Close()

	return content.Bytes(), nil
}

func (export *ODSExport) ContentType() string {
	return "application/vnd.oasis.opendocument.spreadsheet"
}

func odsMetaInfManifestXmlFileContent() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
    <manifest:manifest xmlns:manifest="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0" manifest:version="1.2">
    <manifest:file-entry manifest:full-path="/" manifest:version="1.2" manifest:media-type="application/vnd.oasis.opendocument.spreadsheet"/>
    <manifest:file-entry manifest:full-path="settings.xml" manifest:media-type="text/xml"/>
    <manifest:file-entry manifest:full-path="content.xml" manifest:media-type="text/xml"/>
    <manifest:file-entry manifest:full-path="styles.xml" manifest:media-type="text/xml"/>
    <manifest:file-entry manifest:full-path="meta.xml" manifest:media-type="text/xml"/>
</manifest:manifest>`
}

func odsContentXmlFileContent(e *Experiment) string {
	content := new(bytes.Buffer)

	header := `<?xml version="1.0" encoding="UTF-8"?>
	<office:document-content xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0" xmlns:style="urn:oasis:names:tc:opendocument:xmlns:style:1.0" xmlns:text="urn:oasis:names:tc:opendocument:xmlns:text:1.0" xmlns:table="urn:oasis:names:tc:opendocument:xmlns:table:1.0" xmlns:fo="urn:oasis:names:tc:opendocument:xmlns:xsl-fo-compatible:1.0" xmlns:svg="urn:oasis:names:tc:opendocument:xmlns:svg-compatible:1.0" office:version="1.2">
		<office:scripts/>
		<office:font-face-decls>
			<style:font-face style:name="Arial" svg:font-family="Arial" style:font-family-generic="swiss" style:font-pitch="variable"/>
			<style:font-face style:name="Bitstream Vera Sans" svg:font-family="'Bitstream Vera Sans'" style:font-family-generic="system" style:font-pitch="variable"/>
			<style:font-face style:name="DejaVu Sans" svg:font-family="'DejaVu Sans'" style:font-family-generic="system" style:font-pitch="variable"/>
			<style:font-face style:name="Droid Sans" svg:font-family="'Droid Sans'" style:font-family-generic="system" style:font-pitch="variable"/>
			<style:font-face style:name="FreeSans" svg:font-family="FreeSans" style:font-family-generic="system" style:font-pitch="variable"/>
		</office:font-face-decls>
		<office:automatic-styles>
			<style:style style:name="co1" style:family="table-column">
				<style:table-column-properties fo:break-before="auto" style:column-width="0.889in"/>
			</style:style>
			<style:style style:name="ro1" style:family="table-row">
				<style:table-row-properties style:row-height="0.1681in" fo:break-before="auto" style:use-optimal-row-height="true"/>
			</style:style>
			<style:style style:name="ro2" style:family="table-row">
				<style:table-row-properties style:row-height="0.178in" fo:break-before="auto" style:use-optimal-row-height="true"/>
			</style:style>
			<style:style style:name="ta1" style:family="table" style:master-page-name="Default">
				<style:table-properties table:display="true" style:writing-mode="lr-tb"/>
			</style:style>
		</office:automatic-styles>
		<office:body>
			<office:spreadsheet>
				<table:table table:name="Results" table:style-name="ta1">`
	content.WriteString(header)

	endogenousControlHeader := `<table:table-row table:style-name="ro1">
                    <table:table-cell office:value-type="string">
                        <text:p>name</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>mean</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>stddev</text:p>
                    </table:table-cell>
                    <table:table-cell table:number-columns-repeated="6"/>
                </table:table-row>`
	content.WriteString(endogenousControlHeader)

	for endogenousControlName, endogenousControl := range e.EndogenousControls {
		endogenousControlRow := `<table:table-row table:style-name="ro1">
                    <table:table-cell office:value-type="string">
                        <text:p>%s</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                    <table:table-cell table:number-columns-repeated="6"/>
                </table:table-row>`
		content.WriteString(fmt.Sprintf(endogenousControlRow, endogenousControlName, endogenousControl.Mean, endogenousControl.Mean, endogenousControl.StdDev, endogenousControl.StdDev))
	}

	targetGeneHeader := `<table:table-row table:style-name="ro1">
                    <table:table-cell table:number-columns-repeated="9"/>
                </table:table-row>
                <table:table-row table:style-name="ro1">
                    <table:table-cell office:value-type="string">
                        <text:p>detector</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>name</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>mean</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>stddev</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>dct</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>ddct</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>ddcterr</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>rq</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>rqerr</text:p>
                    </table:table-cell>
                </table:table-row>`
	content.WriteString(targetGeneHeader)

	for detectorName, detector := range e.Detectors {
		for targetGeneName, targetGene := range detector {
			targetGeneRow := `<table:table-row table:style-name="ro1">
                    <table:table-cell office:value-type="string">
                        <text:p>%s</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="string">
                        <text:p>%s</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                    <table:table-cell office:value-type="float" office:value="%f">
                        <text:p>%f</text:p>
                    </table:table-cell>
                </table:table-row>`
			content.WriteString(fmt.Sprintf(targetGeneRow, detectorName, targetGeneName, targetGene.Mean, targetGene.Mean, targetGene.StdDev, targetGene.StdDev, targetGene.DCt, targetGene.DCt, targetGene.DdCt, targetGene.DdCt, targetGene.DdCtErr, targetGene.DdCtErr, targetGene.RQ, targetGene.RQ, targetGene.RQErr, targetGene.RQErr))
		}
	}

	footer := `</table:table>
      <table:named-expressions/>
    </office:spreadsheet>
  </office:body>
</office:document-content>`
	content.WriteString(footer)

	return content.String()
}

func odsMetaXmlFileContent() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<office:document-meta xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:meta="urn:oasis:names:tc:opendocument:xmlns:meta:1.0" office:version="1.2">
    <office:meta>
        <meta:creation-date>2012-11-28T11:33:04</meta:creation-date>
        <dc:date>2012-11-28T11:33:48</dc:date>
        <meta:editing-duration>P0D</meta:editing-duration>
        <meta:editing-cycles>1</meta:editing-cycles>
        <meta:document-statistic meta:table-count="3" meta:cell-count="4" meta:object-count="0"/>
        <meta:generator>AvocadoLab</meta:generator>
    </office:meta>
</office:document-meta>`
}

func odsMimetypeFileContent() string {
	return `application/vnd.oasis.opendocument.spreadsheet`
}

func odsSettingsXmlFileContent() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<office:document-settings xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0" xmlns:config="urn:oasis:names:tc:opendocument:xmlns:config:1.0" xmlns:ooo="http://openoffice.org/2004/office" office:version="1.2">
    <office:settings>
        <config:config-item-set config:name="ooo:view-settings">
            <config:config-item config:name="VisibleAreaTop" config:type="int">0</config:config-item>
            <config:config-item config:name="VisibleAreaLeft" config:type="int">0</config:config-item>
            <config:config-item config:name="VisibleAreaWidth" config:type="int">4515</config:config-item>
            <config:config-item config:name="VisibleAreaHeight" config:type="int">853</config:config-item>
            <config:config-item-map-indexed config:name="Views">
                <config:config-item-map-entry>
                    <config:config-item config:name="ViewId" config:type="string">view1</config:config-item>
                    <config:config-item-map-named config:name="Tables">
                        <config:config-item-map-entry config:name="Sheet1">
                            <config:config-item config:name="CursorPositionX" config:type="int">1</config:config-item>
                            <config:config-item config:name="CursorPositionY" config:type="int">2</config:config-item>
                            <config:config-item config:name="HorizontalSplitMode" config:type="short">0
                            </config:config-item>
                            <config:config-item config:name="VerticalSplitMode" config:type="short">0
                            </config:config-item>
                            <config:config-item config:name="HorizontalSplitPosition" config:type="int">0
                            </config:config-item>
                            <config:config-item config:name="VerticalSplitPosition" config:type="int">0
                            </config:config-item>
                            <config:config-item config:name="ActiveSplitRange" config:type="short">2
                            </config:config-item>
                            <config:config-item config:name="PositionLeft" config:type="int">0</config:config-item>
                            <config:config-item config:name="PositionRight" config:type="int">0</config:config-item>
                            <config:config-item config:name="PositionTop" config:type="int">0</config:config-item>
                            <config:config-item config:name="PositionBottom" config:type="int">0</config:config-item>
                            <config:config-item config:name="ZoomType" config:type="short">0</config:config-item>
                            <config:config-item config:name="ZoomValue" config:type="int">100</config:config-item>
                            <config:config-item config:name="PageViewZoomValue" config:type="int">60
                            </config:config-item>
                            <config:config-item config:name="ShowGrid" config:type="boolean">true</config:config-item>
                        </config:config-item-map-entry>
                    </config:config-item-map-named>
                    <config:config-item config:name="ActiveTable" config:type="string">Sheet1</config:config-item>
                    <config:config-item config:name="HorizontalScrollbarWidth" config:type="int">270
                    </config:config-item>
                    <config:config-item config:name="ZoomType" config:type="short">0</config:config-item>
                    <config:config-item config:name="ZoomValue" config:type="int">100</config:config-item>
                    <config:config-item config:name="PageViewZoomValue" config:type="int">60</config:config-item>
                    <config:config-item config:name="ShowPageBreakPreview" config:type="boolean">false
                    </config:config-item>
                    <config:config-item config:name="ShowZeroValues" config:type="boolean">true</config:config-item>
                    <config:config-item config:name="ShowNotes" config:type="boolean">true</config:config-item>
                    <config:config-item config:name="ShowGrid" config:type="boolean">true</config:config-item>
                    <config:config-item config:name="GridColor" config:type="long">12632256</config:config-item>
                    <config:config-item config:name="ShowPageBreaks" config:type="boolean">true</config:config-item>
                    <config:config-item config:name="HasColumnRowHeaders" config:type="boolean">true
                    </config:config-item>
                    <config:config-item config:name="HasSheetTabs" config:type="boolean">true</config:config-item>
                    <config:config-item config:name="IsOutlineSymbolsSet" config:type="boolean">true
                    </config:config-item>
                    <config:config-item config:name="IsSnapToRaster" config:type="boolean">false</config:config-item>
                    <config:config-item config:name="RasterIsVisible" config:type="boolean">false</config:config-item>
                    <config:config-item config:name="RasterResolutionX" config:type="int">1270</config:config-item>
                    <config:config-item config:name="RasterResolutionY" config:type="int">1270</config:config-item>
                    <config:config-item config:name="RasterSubdivisionX" config:type="int">1</config:config-item>
                    <config:config-item config:name="RasterSubdivisionY" config:type="int">1</config:config-item>
                    <config:config-item config:name="IsRasterAxisSynchronized" config:type="boolean">true
                    </config:config-item>
                </config:config-item-map-entry>
            </config:config-item-map-indexed>
        </config:config-item-set>
        <config:config-item-set config:name="ooo:configuration-settings">
            <config:config-item config:name="ShowNotes" config:type="boolean">true</config:config-item>
            <config:config-item config:name="IsDocumentShared" config:type="boolean">false</config:config-item>
            <config:config-item config:name="AllowPrintJobCancel" config:type="boolean">true</config:config-item>
            <config:config-item config:name="ShowZeroValues" config:type="boolean">true</config:config-item>
            <config:config-item config:name="GridColor" config:type="long">12632256</config:config-item>
            <config:config-item config:name="LoadReadonly" config:type="boolean">false</config:config-item>
            <config:config-item config:name="UpdateFromTemplate" config:type="boolean">true</config:config-item>
            <config:config-item config:name="ShowPageBreaks" config:type="boolean">true</config:config-item>
            <config:config-item config:name="ShowGrid" config:type="boolean">true</config:config-item>
            <config:config-item config:name="SaveVersionOnClose" config:type="boolean">false</config:config-item>
            <config:config-item config:name="IsKernAsianPunctuation" config:type="boolean">false</config:config-item>
            <config:config-item config:name="CharacterCompressionType" config:type="short">0</config:config-item>
            <config:config-item config:name="AutoCalculate" config:type="boolean">true</config:config-item>
            <config:config-item config:name="PrinterName" config:type="string"/>
            <config:config-item config:name="PrinterSetup" config:type="base64Binary"/>
            <config:config-item config:name="IsRasterAxisSynchronized" config:type="boolean">true</config:config-item>
            <config:config-item config:name="IsSnapToRaster" config:type="boolean">false</config:config-item>
            <config:config-item config:name="RasterSubdivisionX" config:type="int">1</config:config-item>
            <config:config-item config:name="RasterResolutionY" config:type="int">1270</config:config-item>
            <config:config-item config:name="RasterResolutionX" config:type="int">1270</config:config-item>
            <config:config-item config:name="RasterIsVisible" config:type="boolean">false</config:config-item>
            <config:config-item config:name="LinkUpdateMode" config:type="short">3</config:config-item>
            <config:config-item config:name="RasterSubdivisionY" config:type="int">1</config:config-item>
            <config:config-item config:name="ApplyUserData" config:type="boolean">true</config:config-item>
            <config:config-item config:name="HasColumnRowHeaders" config:type="boolean">true</config:config-item>
            <config:config-item config:name="HasSheetTabs" config:type="boolean">true</config:config-item>
            <config:config-item config:name="IsOutlineSymbolsSet" config:type="boolean">true</config:config-item>
        </config:config-item-set>
    </office:settings>
</office:document-settings>`
}

func odsStylesXmlFileContent() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<office:document-styles xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0" xmlns:style="urn:oasis:names:tc:opendocument:xmlns:style:1.0" xmlns:text="urn:oasis:names:tc:opendocument:xmlns:text:1.0" xmlns:fo="urn:oasis:names:tc:opendocument:xmlns:xsl-fo-compatible:1.0" xmlns:number="urn:oasis:names:tc:opendocument:xmlns:datastyle:1.0" xmlns:svg="urn:oasis:names:tc:opendocument:xmlns:svg-compatible:1.0" office:version="1.2">
    <office:font-face-decls>
        <style:font-face style:name="Arial" svg:font-family="Arial" style:font-family-generic="swiss" style:font-pitch="variable"/>
        <style:font-face style:name="Bitstream Vera Sans" svg:font-family="&apos;Bitstream Vera Sans&apos;" style:font-family-generic="system" style:font-pitch="variable"/>
        <style:font-face style:name="DejaVu Sans" svg:font-family="&apos;DejaVu Sans&apos;" style:font-family-generic="system" style:font-pitch="variable"/>
        <style:font-face style:name="Droid Sans" svg:font-family="&apos;Droid Sans&apos;" style:font-family-generic="system" style:font-pitch="variable"/>
        <style:font-face style:name="FreeSans" svg:font-family="FreeSans" style:font-family-generic="system" style:font-pitch="variable"/>
    </office:font-face-decls>
    <office:styles>
        <style:default-style style:family="table-cell">
            <style:paragraph-properties style:tab-stop-distance="0.5in"/>
            <style:text-properties style:font-name="Arial" fo:language="en" fo:country="US" style:font-name-asian="Bitstream Vera Sans" style:language-asian="zh" style:country-asian="CN" style:font-name-complex="DejaVu Sans" style:language-complex="hi" style:country-complex="IN"/>
        </style:default-style>
        <number:number-style style:name="N0">
            <number:number number:min-integer-digits="1"/>
        </number:number-style>
        <number:currency-style style:name="N104P0" style:volatile="true">
            <number:currency-symbol number:language="en" number:country="US">$</number:currency-symbol>
            <number:number number:decimal-places="2" number:min-integer-digits="1" number:grouping="true"/>
        </number:currency-style>
        <number:currency-style style:name="N104">
            <style:text-properties fo:color="#ff0000"/>
            <number:text>-</number:text>
            <number:currency-symbol number:language="en" number:country="US">$</number:currency-symbol>
            <number:number number:decimal-places="2" number:min-integer-digits="1" number:grouping="true"/>
            <style:map style:condition="value()&gt;=0" style:apply-style-name="N104P0"/>
        </number:currency-style>
        <style:style style:name="Default" style:family="table-cell">
            <style:text-properties style:font-name-asian="Droid Sans" style:font-name-complex="FreeSans"/>
        </style:style>
        <style:style style:name="Result" style:family="table-cell" style:parent-style-name="Default">
            <style:text-properties fo:font-style="italic" style:text-underline-style="solid" style:text-underline-width="auto" style:text-underline-color="font-color" fo:font-weight="bold"/>
        </style:style>
        <style:style style:name="Result2" style:family="table-cell" style:parent-style-name="Result" style:data-style-name="N104"/>
        <style:style style:name="Heading" style:family="table-cell" style:parent-style-name="Default">
            <style:table-cell-properties style:text-align-source="fix" style:repeat-content="false"/>
            <style:paragraph-properties fo:text-align="center"/>
            <style:text-properties fo:font-size="16pt" fo:font-style="italic" fo:font-weight="bold"/>
        </style:style>
        <style:style style:name="Heading1" style:family="table-cell" style:parent-style-name="Heading">
            <style:table-cell-properties style:rotation-angle="90"/>
        </style:style>
    </office:styles>
    <office:automatic-styles>
        <style:page-layout style:name="Mpm1">
            <style:page-layout-properties style:writing-mode="lr-tb"/>
            <style:header-style>
                <style:header-footer-properties fo:min-height="0.2953in" fo:margin-left="0in" fo:margin-right="0in" fo:margin-bottom="0.0984in"/>
            </style:header-style>
            <style:footer-style>
                <style:header-footer-properties fo:min-height="0.2953in" fo:margin-left="0in" fo:margin-right="0in" fo:margin-top="0.0984in"/>
            </style:footer-style>
        </style:page-layout>
        <style:page-layout style:name="Mpm2">
            <style:page-layout-properties style:writing-mode="lr-tb"/>
            <style:header-style>
                <style:header-footer-properties fo:min-height="0.2953in" fo:margin-left="0in" fo:margin-right="0in" fo:margin-bottom="0.0984in" fo:border="2.49pt solid #000000" fo:padding="0.0071in" fo:background-color="#c0c0c0">
                    <style:background-image/>
                </style:header-footer-properties>
            </style:header-style>
            <style:footer-style>
                <style:header-footer-properties fo:min-height="0.2953in" fo:margin-left="0in" fo:margin-right="0in" fo:margin-top="0.0984in" fo:border="2.49pt solid #000000" fo:padding="0.0071in" fo:background-color="#c0c0c0">
                    <style:background-image/>
                </style:header-footer-properties>
            </style:footer-style>
        </style:page-layout>
    </office:automatic-styles>
    <office:master-styles>
        <style:master-page style:name="Default" style:page-layout-name="Mpm1">
            <style:header>
                <text:p>
                    <text:sheet-name>???</text:sheet-name>
                </text:p>
            </style:header>
            <style:header-left style:display="false"/>
            <style:footer>
                <text:p>Page
                    <text:page-number>1</text:page-number>
                </text:p>
            </style:footer>
            <style:footer-left style:display="false"/>
        </style:master-page>
        <style:master-page style:name="Report" style:page-layout-name="Mpm2">
            <style:header>
                <style:region-left>
                    <text:p>
                        <text:sheet-name>???</text:sheet-name>
                        (<text:title>???</text:title>)
                    </text:p>
                </style:region-left>
                <style:region-right>
                    <text:p><text:date style:data-style-name="N2" text:date-value="2012-11-28">00/00/0000</text:date>,
                        <text:time>00:00:00</text:time>
                    </text:p>
                </style:region-right>
            </style:header>
            <style:header-left style:display="false"/>
            <style:footer>
                <text:p>Page
                    <text:page-number>1</text:page-number>
                    /
                    <text:page-count>99</text:page-count>
                </text:p>
            </style:footer>
            <style:footer-left style:display="false"/>
        </style:master-page>
    </office:master-styles>
</office:document-styles>`
}
