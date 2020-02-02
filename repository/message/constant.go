package message

const STATUS_SUCESS STATUS = "0.DOIP/Status.001"
const STATUS_REQUEST_INVALID STATUS = "0.DOIP/Status.101"
const STATUS_NO_AUTHENTICATE STATUS = "0.DOIP/Status.102"
const STATUS_UNAUTHORIZED STATUS = "0.DOIP/Status.103"
const STATUS_DO_NO_EXIST STATUS = "0.DOIP/Status.104"
const STATUS_DUPLICATED_ID STATUS = "0.DOIP/Status.105"
const STATUS_DECLINED STATUS = "0.DOIP/Status.200"
const STATUS_MULTI_ERRORS STATUS = "0.DOIP/Status.500"

const TYPE_DATA DOType = "data"
const TYPE_CONTRACT DOType = "contract"

const OPERATION_HELLO OPERATION = "0.DOIP/Op.Hello"
const OPERATION_CREATE OPERATION = "0.DOIP/Op.Create"
const OPERATION_RETRIEVE OPERATION = "0.DOIP/Op.Retrieve"
const OPERATION_DELETE OPERATION = "0.DOIP/Op.Delete"
const OPERATION_UPDATE OPERATION = "0.DOIP/Op.Update"
const OPERATION_SEARCH OPERATION = "0.DOIP/Op.Search"
const OPERATION_LISTOPERATIONS OPERATION = "0.DOIP/Op.ListOperations"

const HandleServer string = "http://47.106.38.23:10001/"
const RegisterUrl string = HandleServer + "register/"
const ResolveUrl string = HandleServer + "resolve?"
