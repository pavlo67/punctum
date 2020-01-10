package tagger_pg

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/strlib"
)

const onCountOnTag = "on tagsSQLite.countTag(): "

func (taggerOp *tagsSQLite) countTag(tagLabel string, passedTags []string, labelsRemoved []string, stmPartedSize, stmList, stmAddTag *sql.Stmt) ([]string, error) {
	if strlib.In(passedTags, tagLabel) {
		return passedTags, nil
	}

	var partedSize uint64
	var partedSizePtr *uint64
	values := []interface{}{tagLabel}
	row := stmPartedSize.QueryRow(values...)
	if err := row.Scan(&partedSizePtr); err != nil && err != sql.ErrNoRows {
		return passedTags, errors.Wrapf(err, onCountOnTag+": can't tx.QueryRow(%s, %#v)", taggerOp.sqlTagPartedSize, values)
	}
	if partedSizePtr != nil {
		partedSize = *partedSizePtr
	}

	var labelsOnTag []string
	values = []interface{}{taggerOp.ownInterfaceKey, tagLabel}
	rows, err := stmList.Query(values...)
	if err != sql.ErrNoRows && err != nil {
		return passedTags, errors.Wrapf(err, onCountOnTag+": can't tx.Query(%s, %#v)", taggerOp.sqlList, values)
	}
	defer rows.Close()
	for rows.Next() {
		var tagLabel, relation string
		err = rows.Scan(&tagLabel, &relation)
		if err != nil {
			return passedTags, errors.Wrapf(err, onCountOnTag+": can't tx.ScanQueryRow(%s, %#v)", taggerOp.sqlList, values)
		}
		labelsOnTag = append(labelsOnTag, tagLabel)
	}
	err = rows.Err()
	if err != nil {
		return passedTags, errors.Wrapf(err, onCountOnTag+": "+sqllib.RowsError, taggerOp.sqlList, values)
	}

	values = []interface{}{tagLabel, len(labelsOnTag), partedSize}
	if _, err := stmAddTag.Exec(values...); err != nil {
		return passedTags, errors.Wrapf(err, onCountOnTag+": can't tx.Exec(%s, %#v)", taggerOp.sqlAddTag, values)
	}
	// TODO: don't forget! this must be done before loop with (top!) labelsToCount
	passedTags = append(passedTags, tagLabel)

	labelsToCount := labelsOnTag
	for _, labelRemoved := range labelsRemoved {
		if !strlib.In(labelsToCount, labelRemoved) {
			labelsToCount = append(labelsToCount, labelRemoved)
		}
	}

	for _, labelToCount := range labelsToCount {
		if passedTags, err = taggerOp.countTag(labelToCount, passedTags, nil, stmPartedSize, stmList, stmAddTag); err != nil {
			return passedTags, errors.Wrapf(err, "on tag '%s'", tagLabel)
		}
	}

	return passedTags, nil
}

const onCountChanged = "on tagsSQLite.countTagChanged(): "

func (taggerOp *tagsSQLite) countTagChanged(key joiner.InterfaceKey, id common.ID, tagLabelsRemoved []string, tx *sql.Tx) error {
	if key != taggerOp.ownInterfaceKey {
		return nil
	}
	tagLabel := string(id)

	stmTagParterSize, err := tx.Prepare(taggerOp.sqlTagPartedSize)
	if err != nil {
		return errors.Wrapf(err, onCountChanged+": can't tx.Prepare(%s)", taggerOp.sqlTagPartedSize)
	}
	stmList, err := tx.Prepare(taggerOp.sqlList)
	if err != nil {
		return errors.Wrapf(err, onCountChanged+": can't tx.Prepare(%s)", taggerOp.sqlList)
	}
	stmAddTag, err := tx.Prepare(taggerOp.sqlAddTag)
	if err != nil {
		return errors.Wrapf(err, onCountChanged+": can't tx.Prepare(%s)", taggerOp.sqlAddTag)
	}

	if _, err := taggerOp.countTag(tagLabel, nil, tagLabelsRemoved, stmTagParterSize, stmList, stmAddTag); err != nil {
		return errors.Wrapf(err, onCountChanged+": can't taggerOp.countTag(%s, ...)", tagLabel)
	}

	return nil
}