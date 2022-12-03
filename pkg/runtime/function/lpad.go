/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package function

import (
	"context"
	"fmt"
	"strings"
)

import (
	gxbig "github.com/dubbogo/gost/math/big"

	"github.com/pkg/errors"
)

import (
	"github.com/arana-db/arana/pkg/proto"
	"github.com/arana-db/arana/pkg/util/runes"
)

// FuncLpad is https://dev.mysql.com/doc/refman/5.6/en/string-functions.html#function_lpad
const FuncLpad = "LPAD"

var _ proto.Func = (*lpadFunc)(nil)

func init() {
	proto.RegisterFunc(FuncLpad, lpadFunc{})
}

type lpadFunc struct{}

func (a lpadFunc) Apply(ctx context.Context, inputs ...proto.Valuer) (proto.Value, error) {
	if len(inputs) != 3 {
		return "", errors.New("The Lpad function must accept three parameters\n")
	}

	strInput, err := inputs[0].Value(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lenInput, err := inputs[1].Value(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	padstrInput, err := inputs[2].Value(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	lenNum, err := gxbig.NewDecFromString(fmt.Sprint(lenInput))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if lenNum.IsNegative() {
		return "NULL", nil
	} else if lenNum.IsZero() {
		return "", nil
	}

	if !strings.Contains(fmt.Sprint(lenInput), ".") {
		num, err := lenNum.ToInt()
		if err != nil {
			return "", err
		}
		return a.getLpadStr(runes.ConvertToRune(strInput), num, runes.ConvertToRune(padstrInput))
	}
	num, err := lenNum.ToFloat64()
	if err != nil {
		return "", err
	}
	return a.getLpadStr(runes.ConvertToRune(strInput), int64(num), runes.ConvertToRune(padstrInput))
}

func (a lpadFunc) NumInput() int {
	return 3
}

func (a lpadFunc) getLpadStr(StrRunes []rune, num int64, PadstrRunes []rune) (string, error) {
	if num <= int64(len(StrRunes)) {
		return string(StrRunes[:num]), nil
	} else {
		lenStrRunes := int64(len(StrRunes))
		lenAppend := num - int64(len(StrRunes))
		result := ""
		for lenAppend > int64(len(PadstrRunes)) {
			result = result + string(PadstrRunes[:int64(len(PadstrRunes))])
			lenAppend -= int64(len(PadstrRunes))
		}
		result = result + string(PadstrRunes[:lenAppend]) + string(StrRunes[:lenStrRunes])
		return result, nil
	}
}
