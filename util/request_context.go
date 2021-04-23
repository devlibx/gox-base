package util

type RequestContext interface {
	GetTenant() string
	GetProgram() string
	GetSubProgram() string
	GetCountry() string
	GetCity() string
	GetUdf1() string
	GetUdf2() string
	GetUdf3() string
	GetUdf4() string
	GetUdf5() string
	GetVersion() int
}

type requestContext struct {
	tenant     string
	program    string
	subProgram string
	country    string
	city       string
	udf1       string
	udf2       string
	udf3       string
	udf4       string
	udf5       string
	version    int
}

func (r *requestContext) GetTenant() string {
	return r.tenant
}

func (r *requestContext) GetProgram() string {
	return r.program
}

func (r *requestContext) GetSubProgram() string {
	return r.subProgram
}

func (r *requestContext) GetCountry() string {
	return r.country
}

func (r *requestContext) GetCity() string {
	return r.city
}

func (r *requestContext) GetUdf1() string {
	return r.udf1
}

func (r *requestContext) GetUdf2() string {
	return r.udf2
}

func (r *requestContext) GetUdf3() string {
	return r.udf3
}

func (r *requestContext) GetUdf4() string {
	return r.udf4
}

func (r *requestContext) GetUdf5() string {
	return r.udf5
}

func (r *requestContext) GetVersion() int {
	return r.version
}

type requestContextBuilder struct {
	rc requestContext
}

func (r *requestContextBuilder) Tenant(tenant string) *requestContextBuilder {
	r.rc.tenant = tenant
	return r
}

func (r *requestContextBuilder) Program(program string) *requestContextBuilder {
	r.rc.program = program
	return r
}
func (r *requestContextBuilder) SubProgram(subProgram string) *requestContextBuilder {
	r.rc.subProgram = subProgram
	return r
}
func (r *requestContextBuilder) Country(country string) *requestContextBuilder {
	r.rc.country = country
	return r
}
func (r *requestContextBuilder) City(city string) *requestContextBuilder {
	r.rc.city = city
	return r
}
func (r *requestContextBuilder) Udf1(udf1 string) *requestContextBuilder {
	r.rc.udf1 = udf1
	return r
}
func (r *requestContextBuilder) Udf2(udf2 string) *requestContextBuilder {
	r.rc.udf2 = udf2
	return r
}
func (r *requestContextBuilder) Udf3(udf3 string) *requestContextBuilder {
	r.rc.udf3 = udf3
	return r
}
func (r *requestContextBuilder) Udf4(udf4 string) *requestContextBuilder {
	r.rc.udf4 = udf4
	return r
}
func (r *requestContextBuilder) Udf5(udf5 string) *requestContextBuilder {
	r.rc.udf5 = udf5
	return r
}
func (r *requestContextBuilder) Version(version int) *requestContextBuilder {
	r.rc.version = version
	return r
}

func (r *requestContextBuilder) Build() RequestContext {
	return &r.rc
}

func NewRequestContextBuilder() *requestContextBuilder {
	return &requestContextBuilder{rc: requestContext{
		tenant:     "*",
		program:    "*",
		subProgram: "*",
		country:    "*",
		city:       "*",
		udf1:       "*",
		udf2:       "*",
		udf3:       "*",
		udf4:       "*",
		udf5:       "*",
		version:    1,
	}}
}
